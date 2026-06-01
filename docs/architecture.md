# Architecture And Data Model

## Runtime Shape

The POC runs as two local processes:

- Go API on `http://localhost:8080`
- React frontend on `http://localhost:5173`

Postgres runs through Docker Compose.

## Backend Modules

- `cmd/api`: API entrypoint.
- `cmd/seed`: loads the spreadsheet POC snapshot.
- `internal/config`: environment configuration.
- `internal/domain`: money type, app models, and calculation logic.
- `internal/httpapi`: routes and handlers.
- `internal/store`: persistence interface and Postgres implementation.
- `internal/seed`: structured seed data.

## Database Model

- `households`
- `members`
- `monthly_snapshots`
- `income_entries`
- `budget_categories`
- `budget_items`
- `allocation_rules`
- `household_goals`
- `joint_account_items`
- `joint_account_contributions`

Money is stored as integer pence.

## API Surface

- `GET /api/summary`
- `GET /api/members`
- `GET /api/members/{id}/budget`
- `PUT /api/members/{id}/budget-items/{itemId}`
- `GET /api/joint-account`
- `PUT /api/joint-account/items/{itemId}`
- `GET /api/goals`
- `PUT /api/goals/{id}`

## Frontend Routes

- `/`: household dashboard
- `/members/:memberId`: member budget
- `/joint-account`: shared bills and contributions
- `/goals`: runway and emergency targets

