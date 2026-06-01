# Tech Stack Decision

## Decision

Swaledale uses a Go API, React frontend, and Postgres database for the POC.

## Rationale

- Go is the learning focus for the backend.
- Postgres gives the project a real relational data model from the start.
- React keeps the UI flexible without making the Go app responsible for all interaction state.
- A dev-only POC avoids deployment and auth complexity until the household model is proven.

## Backend

- `chi` for routing.
- `pgx` for Postgres access.
- `goose` for migrations.
- `sqlc` configured for future generated query code.

## Frontend

- React + TypeScript + Vite on Node 24.
- React Router for top-level screens.
- TanStack Query for loading, caching, mutation, and invalidation.
- Tailwind CSS for a quiet dashboard UI.

## Deferred Decisions

- Authentication approach.
- Deployment target.
- Bank import strategy.
- Whether Google Sheets remains an import source after the POC.
