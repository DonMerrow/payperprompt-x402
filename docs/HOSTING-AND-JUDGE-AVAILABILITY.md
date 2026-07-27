# Hosting and Judge Availability

The one-command launcher creates a Cloudflare Quick Tunnel to the local Go
service. It is convenient for a live demo, but it has no uptime guarantee and
stops when the launcher, computer, network, or temporary tunnel stops.

## Submission links that remain useful

Use permanent evidence as the foundation of the submission:

1. public GitHub repository;
2. embedded YouTube demo;
3. static proof page published from `docs/`;
4. verified BaseScan transaction;
5. optional live application URL.

Do not use a temporary `trycloudflare.com` address as the only Try It Out link.
If it is included, label it **live when demo host is online**.

## Publish the static proof page with GitHub Pages

After pushing the repository:

1. open the repository's **Settings → Pages**;
2. choose **Deploy from a branch**;
3. select the public branch and `/docs`;
4. save and wait for the Pages URL;
5. open it in a private browser window;
6. add that stable URL to Devpost.

The page is read-only. It does not run Ollama, open a wallet, or settle a
payment. It preserves the product explanation, architecture, representative
transaction, safety boundary, and reproduction instructions while the
interactive host is offline.

## Stable interactive hosting

The complete application needs long-running Go, Rust, Ollama, and official x402
services. To keep it interactive continuously, use an always-on machine or VPS
and a stable HTTPS endpoint. A named Cloudflare Tunnel gives a stable hostname,
but the origin services must still remain online.

Before exposing an always-on deployment:

- use a disposable testnet payer and merchant;
- protect administrative policy endpoints;
- use managed secrets and per-tenant authentication;
- add rate limiting, logs, monitoring, backups, and incident response;
- complete an external security review;
- do not enable mainnet by merely changing an environment variable.

For this hackathon, the repository, video, static page, and BaseScan proof make
the project reviewable without pretending a laptop-hosted MVP has production
uptime.
