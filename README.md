# Petshop Web-based Management System

Monorepo for a simple, practical pet management platform with appointments, QR public pet profiles, and a lightweight e‑commerce flow.

## Structure

```bash
backend/   # Go API (chi), JWT auth, rate limiting
frontend/  # Vue 3 + TS SPA
database/  # SQL schema (Postgres / Supabase)
docs/      # Diagrams and methodology (Mermaid in Markdown)
```

## Docs

- See `docs/` for:
  - methodology.md
  - context-diagram.md
  - data-flow-diagram-l0.md, data-flow-diagram-l1.md
  - activity-diagram.md
  - entity-relationship-diagram.md
  - database-schema.md

## Quickstart

- Backend

  - env: copy `backend/env.example` to `backend/.env`
  - run: `go run backend/cmd/api/main.go`

- Frontend
  - env: set `VITE_API_URL` to your API base (e.g., `http://localhost:3000`)
  - run: `cd frontend && npm install && npm run dev`

For full features/endpoints, see `backend/README.md` and `PRD.md`.
