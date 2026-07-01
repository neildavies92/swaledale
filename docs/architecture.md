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

## Auth

- Users register with email and password (bcrypt hashed); registration bootstraps a
  household, a member for the user, and an empty starter snapshot.
- Sessions are random tokens stored hashed in `sessions`, delivered as a 30-day
  HttpOnly `swaledale_session` cookie.
- Session middleware resolves the user, and every query and update is scoped to the
  user's household.

## Database Model

- `households`
- `users`
- `sessions`
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

## Frontend Routes

- `/login`: sign in
- `/register`: create an account and household
- `/`: household dashboard (all routes below require a session)
- `/members/:memberId`: member budget
- `/joint-account`: shared bills and contributions
- `/goals`: runway and emergency targets

