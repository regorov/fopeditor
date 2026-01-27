# fopeditor

Developer-focused web editor for building PDF templates with XSL-FO, XML and Apache FOP. The project ships a Go backend, React/Mantine frontend, and a dedicated Apache FOP sidecar image packed for Docker/Kubernetes style deployments.

## Architecture

- **backend** – Go 1.21 REST API with `/health` and `/api/render` endpoints. Rendering flows through `internal/render`, which now calls an HTTP-based FOP sidecar whenever `FOP_ENDPOINT` is configured (falling back to a stub otherwise).
- **frontend** – React + Vite + Mantine single-page editor with dual Monaco editors, render + load-example actions, and inline error handling.
- **fop** – Temurin-based image bundling Apache FOP and a lightweight HTTP wrapper that accepts XML/XSL payloads and streams the rendered PDF back to callers.
- **deploy** – Production `docker-compose` specification wiring the three services on an internal network.
- **CI/CD** – GitHub Actions for testing, publishing container images to GHCR, and remote deployment via SSH.

## Local development

1. **FOP sidecar (optional but recommended)**
   ```bash
   docker build -t fopeditor-fop ./fop
   docker run --rm -p 8090:8090 -v $(pwd)/fonts:/extra-fonts fopeditor-fop
   ```
   The container exposes `POST /render` that the backend can consume at `http://localhost:8090/render`.

2. **Backend**
   ```bash
   cd backend
   go test ./...
   FOP_ENDPOINT=http://localhost:8090/render go run ./cmd/server
   ```
   The server listens on `localhost:8080`. If `FOP_ENDPOINT` is omitted, a stub PDF is generated for rapid iteration.

3. **Frontend**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   Vite serves the UI on `http://localhost:5173` and proxies `/api` calls to the backend.

4. **Rendering**
   With `FOP_ENDPOINT` configured the backend streams documents from the sidecar. Keep it unset for a stubbed PDF that simply reports the payload sizes—helpful when FOP is unavailable locally.

## Docker usage

- **Backend image**
  ```bash
  docker build -t fopeditor-backend ./backend
  docker run -p 8080:8080 fopeditor-backend
  ```
- **Frontend image**
  ```bash
  docker build -t fopeditor-frontend ./frontend
  docker run -p 80:80 fopeditor-frontend
  ```
- **FOP sidecar**
  ```bash
  docker build -t fopeditor-fop ./fop
  docker run --rm -p 8090:8090 -v $(pwd)/fonts:/extra-fonts fopeditor-fop
  ```
  The container exposes `http://localhost:8090/render` for the backend. You can still exec into it and run `fop` manually if desired.
- **Production compose**
  ```bash
  cd deploy
  cat <<'ENV' > .env
  SERVER_NAME=app.example.com
  CERTBOT_EMAIL=admin@example.com
  ENV
  docker compose -f docker-compose.prod.yml up -d
  ```
  The compose file uses the GHCR images that GitHub Actions publishes on pushes to `main`. TLS termination is handled directly inside the frontend container via Let's Encrypt/Certbot, so be sure DNS already points to the host before running the stack.

- **Local compose (builds from source)**
  ```bash
  docker compose up --build
  ```
  This spins up the FOP sidecar, Go backend, and frontend dev proxy in one command. The UI is available at `http://localhost:5173`, while the API listens on `http://localhost:8080`.

## TLS and Let's Encrypt

The frontend container can automatically provision and renew Let's Encrypt certificates.

- Set `ENABLE_TLS=true`, `SERVER_NAME=your.domain`, and `CERTBOT_EMAIL=you@example.com` on the frontend service (already wired in `deploy/docker-compose.prod.yml` via `.env` variables).
- Two named volumes (`frontend-certbot-etc`, `frontend-certbot-www`) persist `/etc/letsencrypt` and the ACME challenge directory so renewals survive container restarts.
- On startup the entrypoint requests a certificate via the ACME HTTP challenge (`/.well-known/acme-challenge/*`) and then starts Nginx listening on ports 80/443.
- A lightweight renewal loop runs every 12 hours (tunable via `CERTBOT_RENEW_INTERVAL_HOURS`). Successful renewals trigger `nginx -s reload` to pick up the new certificates without downtime.

For environments where TLS termination is handled elsewhere (local development, reverse proxy, load balancer), simply omit `ENABLE_TLS` (default is HTTP-only mode using the original `nginx.conf`).

## Continuous delivery

1. `ci.yml` runs Go tests and a frontend build on pushes and pull requests.
2. `release.yml` builds and pushes backend, frontend, and FOP images to `ghcr.io/regorov`.
3. `deploy.yml` waits for a successful release run, then SSHes into the DigitalOcean host and executes `docker compose pull && docker compose up -d` in `/opt/fopeditor`.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for pull request guidelines. Issues and improvements are welcome!
