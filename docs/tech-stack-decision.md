# Tech Stack Decision

## Decision

Swaledale uses a Go API, React frontend, and Postgres database for household FIRE planning.

## Rationale

- Go is the learning focus for the backend.
- Postgres gives the project a real relational data model from the start.
- React keeps the UI flexible without making the Go app responsible for all interaction state.
- Preserve the working local stack and household-scoped authentication while evolving the domain.

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

- Grafana integration (FIRE-002/FIRE-006); React remains available for management.
- Deployment target.
- Bank import strategy.
- Google Sheets collection is tracked in FIRE-007.
