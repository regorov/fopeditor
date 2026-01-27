# Task: Create OSS project "fopeditor"

You are working in an empty GitHub repository:

    https://github.com/regorov/fopeditor

The project is **open source** and should be fully bootstrapped:
directories, source code, Docker files, CI/CD pipelines, and documentation.

---

## 1. Project goal

Create a web service that allows users to:

- edit **XSL (XSLT 1.0 / XSL-FO)** in one editor
- edit **XML** in another editor
- press **Render PDF**
- receive a rendered **PDF document (Apache FOP)** opened in a new browser tab

The service is intended as a **developer tool / editor**, not a production document storage system.

---

## 2. Technology stack (mandatory)

### Backend
- Go (golang)
- Stateless HTTP API
- REST
- Runs in Docker

### Frontend
- React
- TypeScript
- Vite
- Mantine UI
- Code editor with syntax highlighting (Monaco or CodeMirror)

### PDF rendering
- Apache FOP
- XSLT 1.0 compatible
- Implemented as a **sidecar container** (preferred)

### Infrastructure
- Docker
- docker-compose (local + prod)
- GitHub Actions (CI/CD)
- GitHub Container Registry (ghcr.io)
- Target hosting: DigitalOcean droplet

---

## 3. Repository structure (must be created)

Create the following structure:

fopeditor/
backend/
cmd/server/main.go
internal/http/
internal/render/
go.mod
go.sum
Dockerfile

frontend/
src/
main.tsx
App.tsx
components/Editors.tsx
index.html
package.json
tsconfig.json
vite.config.ts
nginx.conf
Dockerfile

fop/
Dockerfile

deploy/
docker-compose.prod.yml

.github/
workflows/
ci.yml
release.yml
deploy.yml

.editorconfig
.gitignore
LICENSE
README.md
CONTRIBUTING.md

All directories and files above **must be created**.

---

## 4. Backend requirements

### API

Implement HTTP server listening on port `8080`.

#### Healthcheck

GET /health

Response:
- 200 OK
- body: `{ "status": "ok" }`

#### Render PDF

POST /api/render

Request (JSON):
```json
{
  "xsl": "<xsl>...</xsl>",
  "xml": "<xml>...</xml>"
}

Response (success):
	•	HTTP 200
	•	Content-Type: application/pdf
	•	binary PDF body

Response (error):
	•	HTTP 4xx/5xx
	•	JSON:

{
  "code": "XML_ERROR | XSL_ERROR | FOP_ERROR",
  "message": "human readable message"
}

Implementation notes
	•	Backend is stateless
	•	No persistence
	•	Temporary files allowed, must be cleaned up
	•	Rendering logic should be isolated in internal/render
	•	FOP is expected to be available as an external container or CLI (stub allowed for now)

⸻

5. Frontend requirements

UI

Single-page application:
	•	Two side-by-side editors:
	•	left: XSL
	•	right: XML
	•	Button: Render PDF
	•	Button: Load example
	•	Error panel for render errors

Editor
	•	Syntax highlighting for XML/XSL
	•	Line numbers
	•	Reasonable defaults

Behavior
	•	Clicking Render PDF:
	•	sends POST /api/render
	•	opens returned PDF in a new browser tab
	•	Errors are shown in UI

⸻

6. Docker requirements

Backend
	•	Multi-stage Dockerfile
	•	Final image should be minimal (distroless or alpine)
	•	Expose port 8080

Frontend
	•	Build with Node
	•	Serve with Nginx
	•	Proxy /api/* to backend

FOP
	•	Separate image
	•	Based on JRE (Temurin or similar)
	•	Apache FOP installed
	•	No HTTP API required yet (CLI container is enough)

⸻

7. docker-compose (production)

Create deploy/docker-compose.prod.yml with services:
	•	frontend
	•	backend
	•	fop

Services must be connected via internal Docker network.

⸻

8. CI/CD (GitHub Actions)

CI (ci.yml)

Triggered on:
	•	pull_request
	•	push to main

Steps:
	•	backend: go test ./...
	•	frontend: install + build
	•	fail on errors

Build & Push (release.yml)

Triggered on:
	•	push to main

Steps:
	•	build Docker images:
	•	fopeditor-backend
	•	fopeditor-frontend
	•	fopeditor-fop
	•	push images to:

ghcr.io/regorov/*



Deploy (deploy.yml)

Triggered after successful image build.

Steps:
	•	SSH into DigitalOcean droplet
	•	docker compose pull
	•	docker compose up -d

Secrets are assumed to exist:
	•	DO_HOST
	•	DO_USER
	•	DO_SSH_KEY
	•	GHCR_PAT

⸻

9. Open source requirements
	•	Add MIT or Apache-2.0 license
	•	README.md must include:
	•	project description
	•	architecture overview
	•	local development instructions
	•	Docker usage
	•	CONTRIBUTING.md with basic PR rules

⸻

10. Output expectations

After completing this task, the repository should:
	•	Build successfully in CI
	•	Produce Docker images
	•	Be deployable to DigitalOcean
	•	Show a working UI (even with stub PDF rendering initially)

You may stub FOP rendering logic if necessary, but structure and interfaces must be correct.
