const $ = (id) => document.getElementById(id);

async function json(path) {
  const response = await fetch(path, { cache: "no-store" });
  const payload = await response.json();
  return { response, payload };
}

async function load() {
  try {
    const { payload: health } = await json("/api/health");
    $("mode").textContent = health.mode;
    $("network").textContent = health.network;
    $("ai").textContent = health.ollama_ready ? `${health.ai_model} ready` : `${health.ai_model} unavailable`;
    $("price").textContent = `$${health.price_usd} USDC`;
    $("merchant").textContent = health.merchant;

    const { response, payload: proof } = await json("/api/proof");
    if (!response.ok) return;

    $("proof-title").textContent = "Official settlement proof generated";
    $("proof-copy").textContent = "Open BaseScan and confirm the transfer before marking the submission verified.";
    $("explorer").href = proof.explorer_url;
    $("explorer").classList.remove("hidden");
    $("proof-fields").classList.remove("hidden");
    $("payer").textContent = proof.payer;
    $("proof-merchant").textContent = proof.merchant;
    $("transaction").textContent = proof.settlement.transaction;
    $("proof-model").textContent = proof.paid_api_response?.ai_model ?? "unknown";
  } catch (error) {
    $("error").textContent = `Status error: ${error.message}`;
  }
}

load();
