# Contributing to fopeditor

## Getting started
- Fork the repository and clone it locally.
- Install Go 1.21+ and Node.js 20+.
- (Optional) Run the FOP sidecar locally (`docker run --rm -p 8090:8090 fopeditor-fop`) and export `FOP_ENDPOINT=http://localhost:8090/render` to exercise real rendering.
- Run `go test ./...` inside `backend/` and `npm run build` inside `frontend/` before opening a PR.

## Pull requests
- Use feature branches and keep commits focused.
- Update README/documentation when behavior or workflows change.
- Include tests where practical (unit tests for backend packages, UI tests or stories for complex frontend changes).
- Describe what the change does and how to validate it in the PR body.

## Code style
- Follow Go formatting via `gofmt` and lint-friendly patterns.
- Prefer functional React components with hooks and Mantine primitives.
- Keep Dockerfiles minimal and multi-stage whenever possible.

## Communication
- Open an issue for large features or architectural changes before starting work.
- Be respectful and inclusive in all discussions.
