use hmac::{Hmac, Mac};
use serde::{Deserialize, Serialize};
use serde_json::Value;
use sha2::Sha256;
use std::env;
use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};

type HmacSha256 = Hmac<Sha256>;

#[derive(Debug, Clone, Deserialize, Serialize)]
struct Routing {
    route_id: String,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
struct Receipt {
    receipt_id: String,
    request_id: String,
    network: String,
    asset: String,
    amount_usd: String,
    transaction_id: String,
    settled: bool,
    payer: String,
    replay_protected: bool,
    routing: Routing,
    issued_at: String,
    integrity_hmac_sha256: String,
}

#[derive(Debug, Deserialize)]
struct VerifyRequest {
    receipt: Receipt,
}

#[derive(Debug, Serialize)]
struct VerifyResponse {
    valid: bool,
    integrity_valid: bool,
    semantic_valid: bool,
    verifier: &'static str,
    algorithm: &'static str,
    canonical_fields: [&'static str; 9],
    reason: String,
}

#[derive(Debug, Serialize)]
struct OfficialProofResponse {
    valid: bool,
    verifier: &'static str,
    verification_scope: &'static str,
    checks: Vec<ProofCheck>,
    reason: String,
}

#[derive(Debug, Serialize)]
struct ProofCheck {
    name: &'static str,
    passed: bool,
}

fn main() {
    let port = env_value("PAYPERPROMPT_RUST_VERIFY_PORT", "8085");
    let secret = env_value(
        "RECEIPT_HMAC_SECRET",
        "local-dev-receipt-secret-change-me",
    );
    let listener = TcpListener::bind(format!("127.0.0.1:{port}"))
        .expect("bind PayPerPrompt Rust receipt verifier");

    println!("PayPerPrompt Rust receipt verifier listening at http://127.0.0.1:{port}");
    println!("Algorithm: HMAC-SHA256 over 9 canonical receipt fields");

    for stream in listener.incoming() {
        match stream {
            Ok(stream) => handle_connection(stream, &secret),
            Err(error) => eprintln!("connection error: {error}"),
        }
    }
}

fn handle_connection(mut stream: TcpStream, secret: &str) {
    let request = match read_http_request(&mut stream) {
        Ok(request) => request,
        Err(error) => {
            write_json(
                &mut stream,
                400,
                &serde_json::json!({"error": error}).to_string(),
            );
            return;
        }
    };
    let (head, body) = split_request(&request);
    let request_line = head.lines().next().unwrap_or_default();

    if request_line.starts_with("GET /api/health ") {
        write_json(
            &mut stream,
            200,
            &serde_json::json!({
                "ok": true,
                "service": "PayPerPrompt independent Rust receipt verifier",
                "runtime": "rust",
                "algorithm": "HMAC-SHA256",
                "port": env_value("PAYPERPROMPT_RUST_VERIFY_PORT", "8085")
            })
            .to_string(),
        );
        return;
    }

    if request_line.starts_with("POST /api/verify-receipt ") {
        let input: VerifyRequest = match serde_json::from_str(body) {
            Ok(input) => input,
            Err(error) => {
                write_json(
                    &mut stream,
                    400,
                    &serde_json::json!({
                        "valid": false,
                        "error": format!("invalid receipt JSON: {error}")
                    })
                    .to_string(),
                );
                return;
            }
        };
        let result = verify_receipt(&input.receipt, secret);
        let status = if result.valid { 200 } else { 422 };
        write_json(
            &mut stream,
            status,
            &serde_json::to_string_pretty(&result).expect("serialize verification"),
        );
        return;
    }

    if request_line.starts_with("POST /api/verify-official-proof ") {
        let proof: Value = match serde_json::from_str(body) {
            Ok(input) => input,
            Err(error) => {
                write_json(
                    &mut stream,
                    400,
                    &serde_json::json!({
                        "valid": false,
                        "error": format!("invalid official proof JSON: {error}")
                    })
                    .to_string(),
                );
                return;
            }
        };
        let result = verify_official_proof(&proof);
        let status = if result.valid { 200 } else { 422 };
        write_json(
            &mut stream,
            status,
            &serde_json::to_string_pretty(&result).expect("serialize verification"),
        );
        return;
    }

    write_json(
        &mut stream,
        404,
        &serde_json::json!({"error": "not found"}).to_string(),
    );
}

fn verify_official_proof(proof: &Value) -> OfficialProofResponse {
    const NETWORK: &str = "eip155:84532";
    const USDC: &str = "0x036CbD53842c5426634e7929541eC2318f3dCF7e";

    let top_network = text(proof, "/network");
    let top_asset = text(proof, "/asset");
    let payer = text(proof, "/payer");
    let merchant = text(proof, "/merchant");
    let amount = text(proof, "/amount_atomic");
    let transaction = text(proof, "/settlement/transaction");
    let selected_price = text(proof, "/agent_plan/selected/price_usd");
    let expected_atomic = usd_to_atomic(selected_price);

    let checks = vec![
        ProofCheck {
            name: "official proof version",
            passed: text(proof, "/proof_version") == "payperprompt-official-v1",
        },
        ProofCheck {
            name: "Base Sepolia network consistency",
            passed: top_network == NETWORK
                && text(proof, "/settlement/network") == NETWORK
                && text(proof, "/payment_requirement/network") == NETWORK,
        },
        ProofCheck {
            name: "official Base Sepolia USDC contract",
            passed: top_asset.eq_ignore_ascii_case(USDC)
                && text(proof, "/payment_requirement/asset").eq_ignore_ascii_case(USDC),
        },
        ProofCheck {
            name: "distinct valid EVM payer and merchant",
            passed: valid_evm_address(payer)
                && valid_evm_address(merchant)
                && !payer.eq_ignore_ascii_case(merchant),
        },
        ProofCheck {
            name: "settlement identity consistency",
            passed: proof.pointer("/settlement/success").and_then(Value::as_bool) == Some(true)
                && text(proof, "/settlement/payer").eq_ignore_ascii_case(payer)
                && valid_transaction(transaction),
        },
        ProofCheck {
            name: "challenge amount and merchant consistency",
            passed: !amount.is_empty()
                && text(proof, "/payment_requirement/amount") == amount
                && text(proof, "/payment_requirement/payTo").eq_ignore_ascii_case(merchant),
        },
        ProofCheck {
            name: "AI-selected price equals atomic payment",
            passed: expected_atomic.is_some_and(|value| value == amount),
        },
        ProofCheck {
            name: "selected route equals paid API response",
            passed: text(proof, "/agent_plan/selected/route_id")
                == text(proof, "/paid_api_response/route_id")
                && text(proof, "/agent_plan/selected/provider")
                    == text(proof, "/paid_api_response/provider")
                && text(proof, "/agent_plan/selected/path")
                    == text(proof, "/paid_api_response/paid_resource")
                && selected_price == text(proof, "/paid_api_response/price_usd"),
        },
        ProofCheck {
            name: "paid response used AI",
            passed: proof
                .pointer("/paid_api_response/ai_used")
                .and_then(Value::as_bool)
                == Some(true),
        },
        ProofCheck {
            name: "live-chain evidence flags",
            passed: [
                "/live_chain_verification/checked_live",
                "/live_chain_verification/valid",
                "/live_chain_verification/transaction_success",
                "/live_chain_verification/usdc_transfer_matched",
            ]
            .iter()
            .all(|path| proof.pointer(path).and_then(Value::as_bool) == Some(true)),
        },
    ];
    let valid = checks.iter().all(|check| check.passed);
    OfficialProofResponse {
        valid,
        verifier: "payperprompt-official-proof-verifier-rust",
        verification_scope: "independent semantic and cross-field verification of recorded x402 proof; Go performs the live JSON-RPC chain query",
        reason: if valid {
            "Official proof fields, AI route, x402 challenge, settlement, and recorded live-chain evidence are internally consistent."
        } else {
            "One or more official proof consistency checks failed."
        }
        .to_string(),
        checks,
    }
}

fn text<'a>(value: &'a Value, pointer: &str) -> &'a str {
    value.pointer(pointer).and_then(Value::as_str).unwrap_or("")
}

fn valid_evm_address(value: &str) -> bool {
    value.len() == 42
        && value.starts_with("0x")
        && value[2..].chars().all(|character| character.is_ascii_hexdigit())
}

fn valid_transaction(value: &str) -> bool {
    value.len() == 66
        && value.starts_with("0x")
        && value[2..].chars().all(|character| character.is_ascii_hexdigit())
}

fn usd_to_atomic(value: &str) -> Option<String> {
    let (whole, fraction) = value.split_once('.').unwrap_or((value, ""));
    if whole.is_empty()
        || !whole.chars().all(|character| character.is_ascii_digit())
        || !fraction.chars().all(|character| character.is_ascii_digit())
        || fraction.len() > 6
    {
        return None;
    }
    let padded = format!("{fraction:0<6}");
    let atomic = format!("{whole}{padded}").parse::<u128>().ok()?;
    Some(atomic.to_string())
}

fn verify_receipt(receipt: &Receipt, secret: &str) -> VerifyResponse {
    let semantic_valid = receipt.settled
        && receipt.replay_protected
        && !receipt.receipt_id.is_empty()
        && !receipt.request_id.is_empty()
        && !receipt.transaction_id.is_empty()
        && !receipt.payer.is_empty()
        && !receipt.routing.route_id.is_empty()
        && receipt.amount_usd.parse::<f64>().is_ok_and(|amount| amount > 0.0);

    let supplied_tag = decode_hex(&receipt.integrity_hmac_sha256);
    let integrity_valid = supplied_tag.is_some_and(|tag| {
        let mut mac = HmacSha256::new_from_slice(secret.as_bytes())
            .expect("HMAC accepts keys of any length");
        mac.update(canonical_receipt(receipt).as_bytes());
        mac.verify_slice(&tag).is_ok()
    });

    let valid = semantic_valid && integrity_valid;
    let reason = match (semantic_valid, integrity_valid) {
        (true, true) => "Receipt fields and HMAC-SHA256 integrity tag are valid.",
        (false, true) => "Integrity tag matches, but required settlement semantics failed.",
        (true, false) => "Receipt fields are complete, but the integrity tag does not match.",
        (false, false) => "Receipt semantics and integrity tag are invalid.",
    }
    .to_string();

    VerifyResponse {
        valid,
        integrity_valid,
        semantic_valid,
        verifier: "payperprompt-receipt-verifier-rust",
        algorithm: "HMAC-SHA256",
        canonical_fields: [
            "receipt_id",
            "request_id",
            "network",
            "asset",
            "amount_usd",
            "transaction_id",
            "payer",
            "routing.route_id",
            "issued_at",
        ],
        reason,
    }
}

fn canonical_receipt(receipt: &Receipt) -> String {
    [
        receipt.receipt_id.as_str(),
        receipt.request_id.as_str(),
        receipt.network.as_str(),
        receipt.asset.as_str(),
        receipt.amount_usd.as_str(),
        receipt.transaction_id.as_str(),
        receipt.payer.as_str(),
        receipt.routing.route_id.as_str(),
        receipt.issued_at.as_str(),
    ]
    .join("|")
}

fn read_http_request(stream: &mut TcpStream) -> Result<String, String> {
    let mut buffer = Vec::with_capacity(8192);
    let mut chunk = [0_u8; 4096];
    let mut expected_length = None;

    loop {
        let read = stream.read(&mut chunk).map_err(|error| error.to_string())?;
        if read == 0 {
            break;
        }
        buffer.extend_from_slice(&chunk[..read]);
        if buffer.len() > 1_048_576 {
            return Err("request exceeds 1 MiB".to_string());
        }

        if let Some(header_end) = find_bytes(&buffer, b"\r\n\r\n") {
            if expected_length.is_none() {
                let head = String::from_utf8_lossy(&buffer[..header_end]);
                expected_length = Some(content_length(&head));
            }
            let body_length = buffer.len().saturating_sub(header_end + 4);
            if body_length >= expected_length.unwrap_or(0) {
                break;
            }
        }
    }

    String::from_utf8(buffer).map_err(|error| error.to_string())
}

fn content_length(head: &str) -> usize {
    head.lines()
        .filter_map(|line| line.split_once(':'))
        .find(|(name, _)| name.eq_ignore_ascii_case("content-length"))
        .and_then(|(_, value)| value.trim().parse().ok())
        .unwrap_or(0)
}

fn find_bytes(haystack: &[u8], needle: &[u8]) -> Option<usize> {
    haystack
        .windows(needle.len())
        .position(|window| window == needle)
}

fn split_request(request: &str) -> (&str, &str) {
    request.split_once("\r\n\r\n").unwrap_or((request, ""))
}

fn write_json(stream: &mut TcpStream, status: u16, body: &str) {
    let reason = match status {
        200 => "OK",
        400 => "Bad Request",
        404 => "Not Found",
        422 => "Unprocessable Entity",
        _ => "OK",
    };
    let response = format!(
        "HTTP/1.1 {status} {reason}\r\nContent-Type: application/json; charset=utf-8\r\nContent-Length: {}\r\nConnection: close\r\nX-Content-Type-Options: nosniff\r\n\r\n{body}",
        body.len()
    );
    let _ = stream.write_all(response.as_bytes());
}

fn decode_hex(value: &str) -> Option<Vec<u8>> {
    if !value.len().is_multiple_of(2) {
        return None;
    }
    (0..value.len())
        .step_by(2)
        .map(|index| u8::from_str_radix(&value[index..index + 2], 16).ok())
        .collect()
}

fn env_value(key: &str, fallback: &str) -> String {
    env::var(key).unwrap_or_else(|_| fallback.to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    fn sample_receipt() -> Receipt {
        Receipt {
            receipt_id: "receipt-1".to_string(),
            request_id: "request-1".to_string(),
            network: "local-x402-go".to_string(),
            asset: "USDC_TEST".to_string(),
            amount_usd: "0.04".to_string(),
            transaction_id: "go-sandbox-tx-1".to_string(),
            settled: true,
            payer: "sandbox-agent-001".to_string(),
            replay_protected: true,
            routing: Routing {
                route_id: "guardrail-deep".to_string(),
            },
            issued_at: "2026-07-24T01:00:00Z".to_string(),
            integrity_hmac_sha256: String::new(),
        }
    }

    fn sign(receipt: &Receipt, secret: &str) -> String {
        let mut mac =
            HmacSha256::new_from_slice(secret.as_bytes()).expect("valid HMAC key");
        mac.update(canonical_receipt(receipt).as_bytes());
        encode_hex(&mac.finalize().into_bytes())
    }

    fn encode_hex(bytes: &[u8]) -> String {
        bytes.iter().map(|byte| format!("{byte:02x}")).collect()
    }

    #[test]
    fn accepts_valid_receipt() {
        let secret = "test-secret";
        let mut receipt = sample_receipt();
        receipt.integrity_hmac_sha256 = sign(&receipt, secret);
        assert!(verify_receipt(&receipt, secret).valid);
    }

    #[test]
    fn rejects_tampered_amount() {
        let secret = "test-secret";
        let mut receipt = sample_receipt();
        receipt.integrity_hmac_sha256 = sign(&receipt, secret);
        receipt.amount_usd = "400.00".to_string();
        let result = verify_receipt(&receipt, secret);
        assert!(!result.valid);
        assert!(!result.integrity_valid);
    }

    #[test]
    fn rejects_unsettled_receipt() {
        let secret = "test-secret";
        let mut receipt = sample_receipt();
        receipt.integrity_hmac_sha256 = sign(&receipt, secret);
        receipt.settled = false;
        let result = verify_receipt(&receipt, secret);
        assert!(!result.valid);
        assert!(!result.semantic_valid);
    }

    fn official_proof() -> Value {
        serde_json::json!({
            "proof_version": "payperprompt-official-v1",
            "network": "eip155:84532",
            "asset": "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
            "amount_atomic": "20000",
            "payer": "0x826154a3d58aeA3FBD2aa64aAD424594ade927eF",
            "merchant": "0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630",
            "settlement": {
                "success": true,
                "network": "eip155:84532",
                "payer": "0x826154a3d58aeA3FBD2aa64aAD424594ade927eF",
                "transaction": "0x45c95f934bed4ae5dac90f8ba32b28aa1eaa52e66d793f5241a733ce47968c9b"
            },
            "payment_requirement": {
                "network": "eip155:84532",
                "asset": "0x036CbD53842c5426634e7929541eC2318f3dCF7e",
                "amount": "20000",
                "payTo": "0x07fB6cDd24cF265f8ea01A323708DB34d6Dbb630"
            },
            "agent_plan": {"selected": {
                "route_id": "guardrail-fast",
                "provider": "Rapid Policy",
                "path": "/api/services/rapid-policy/check-prompt",
                "price_usd": "0.02"
            }},
            "paid_api_response": {
                "ai_used": true,
                "route_id": "guardrail-fast",
                "provider": "Rapid Policy",
                "paid_resource": "/api/services/rapid-policy/check-prompt",
                "price_usd": "0.02"
            },
            "live_chain_verification": {
                "checked_live": true,
                "valid": true,
                "transaction_success": true,
                "usdc_transfer_matched": true
            }
        })
    }

    #[test]
    fn accepts_consistent_official_proof() {
        assert!(verify_official_proof(&official_proof()).valid);
    }

    #[test]
    fn rejects_tampered_official_amount() {
        let mut proof = official_proof();
        proof["amount_atomic"] = Value::String("40000".to_string());
        assert!(!verify_official_proof(&proof).valid);
    }
}
