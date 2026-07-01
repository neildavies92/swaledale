# Swaledale

Swaledale is a dev-only POC that turns the household bills spreadsheet into a family finance application.

The first release focuses on spreadsheet parity:

- member budgets for Neil and Katie
- shared household bills in the joint account
- emergency fund and runway goals
- editable recurring bill/saving amounts
- Postgres-backed data seeded from the original Google Sheet snapshot
- user registration and sign-in with per-household data isolation

## Stack

- Go API using `chi` and `pgx`
- Postgres via Docker Compose
- React + TypeScript + Vite frontend
- TanStack Query for API state
- Tailwind CSS for the dashboard UI
- `goose` migrations and `sqlc` query definitions

## Prerequisites

- Go `1.24.2` or newer
- Docker
- Postgres client tools are useful but not required
- Node `24.16.0` via `nvm use` for the frontend

## Local Setup

```sh
cp .env.example .env
make db-up
make migrate
make seed
make api
```

In another shell:

```sh
nvm use
cd frontend
npm install
npm run dev
```

Then open `http://localhost:5173`.

## Accounts

All budget data is scoped to a household. Registering creates a fresh household with
an empty starter budget; signing in with a seeded dev account opens the spreadsheet
snapshot household:

- `neil@swaledale.local` / `swaledale-dev`
- `katie@swaledale.local` / `swaledale-dev`

Sessions are 30-day HttpOnly cookies. Set `SECURE_COOKIES=true` when serving over
HTTPS.

## Useful Commands

```sh
make db-up      # start Postgres
make migrate    # run database migrations
make seed       # load the spreadsheet POC snapshot
make api        # run the Go API on :8080
make frontend   # run the Vite dev server
make test       # run backend and frontend tests
make lint       # run backend tests and frontend lint
```

## API

- `GET /healthz`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/auth/me`
- `GET /api/summary`
- `GET /api/members`
- `GET /api/members/{id}/budget`
- `PUT /api/members/{id}/budget-items/{itemId}`
- `GET /api/joint-account`
- `PUT /api/joint-account/items/{itemId}`
- `GET /api/goals`
- `PUT /api/goals/{id}`

All `/api` routes other than register/login require a session cookie.

## Project Docs

See [`docs/project-brief.md`](docs/project-brief.md), [`docs/architecture.md`](docs/architecture.md), [`docs/spreadsheet-data-audit.md`](docs/spreadsheet-data-audit.md), and [`docs/jira-backlog.md`](docs/jira-backlog.md).
