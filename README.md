# Swaledale

Swaledale is a provider-independent household financial independence and FIRE planning platform. The initial deployment serves one household with multiple members, personal finances and jointly owned accounts.

The canonical Go domain now models explicit ownership, stable accounts, replaceable provider mappings, investment instruments, historical balances/holdings, transactions and derived FIRE snapshots. Providers and funds are configuration, not account identity.

The [FIRE Dashboard Project tab](https://docs.google.com/spreadsheets/d/1pStlDjjlz6qp1908SJvJVThnQWj3SH0DzKO70mDza4g/edit?gid=944913493#gid=944913493) is the authoritative backlog and architecture plan.

FIRE-001 defines and validates the domain. FIRE persistence (FIRE-003), configuration (FIRE-024), collectors, calculations and Grafana remain future work. The running React/API screens still serve the earlier budget workflow through `internal/legacybudget`; they are not yet a FIRE dashboard. React is retained for future household/account/connector management.

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

## Development sign-in

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
make seed       # reset local data and load the legacy development fixture
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

See the [project brief](docs/project-brief.md), [canonical architecture and persistence boundary](docs/architecture.md), and [decision log](docs/decision-log.md). The [spreadsheet audit](docs/spreadsheet-data-audit.md) is historical reference only.
