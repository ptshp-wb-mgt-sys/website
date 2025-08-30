# Methodology

This thesis project follows a pragmatic, lightweight methodology focused on shipping a clear, working system for pet management and e‑commerce.

## Goals

- Deliver a small but complete system: backend API, frontend app, and database.
- Keep scope tight and diagrams simple so they’re easy to understand.

## Process

1. Requirements sketch
   - Identify core actors: Client (pet owner), Veterinarian, Public Visitor (QR scan).
   - Identify core features: Pets, QR public profile, Medical Records, Appointments, Products, Orders.
2. Architecture first
   - Define a RESTful API (Go + chi) and a SPA frontend (Vue + Tailwind).
   - Use Supabase Postgres for persistence. Authorization enforced in the backend.
3. Data modeling
   - Start with ERD from core entities and relationships.
   - Translate ERD to SQL tables and indexes.
4. API surface
   - Map features to endpoints under `/api/v1` with JWT protection (except public QR).
   - Keep handlers small, rely on store interfaces for data access.
5. Frontend flows
   - Simple views for Pets, Profiles, Appointments, Products, Orders.
   - Use a thin `apiFetch` wrapper for HTTP calls.
6. Security & reliability
   - JWT auth middleware, rate limiting, CORS, timeouts, and request logging.
   - CRUD validation at handler + store layers.
7. Iteration & testing
   - Manual test scripts and lean handler tests.

## Stack

- Backend: Go, chi router, structured handlers, store interface.
- Frontend: Vue 3 + TypeScript, Tailwind CSS, Pinia, Vue Router.
- Database: Postgres (Supabase), SQL schema stored in `database/schema/tables.sql`.
- Infra: Config‑driven, deployable via Fly.io (sample `fly.toml`).

## Deliverables

- Context diagram, DFD Level 0/1, Activity diagram, ER diagram.
- Database schema summary with references to the full DDL.
