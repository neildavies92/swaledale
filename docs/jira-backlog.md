# Jira-Ready Backlog

## Epic

`Swaledale Household Finance MVP`

Transform the bills spreadsheet into a Go + React household finance application.

## Stories

### Set up Go/Postgres/React project foundation

Acceptance criteria:

- Repo has backend and frontend projects.
- Docker Compose starts Postgres locally.
- Migrations can create the schema.
- README documents local startup.

Tasks:

- Scaffold Go module, API server, config loading, and health endpoint.
- Add Docker Compose Postgres.
- Add migration workflow.
- Scaffold React/Vite frontend.

### Seed spreadsheet POC data

Acceptance criteria:

- Neil, Katie, and Joint Account data loads into Postgres.
- Seeded totals reproduce the spreadsheet snapshot.

Tasks:

- Define relational schema.
- Encode spreadsheet data as a fixed seed snapshot.
- Add seed command.
- Add calculation tests for totals.

### Create household summary dashboard

Acceptance criteria:

- Dashboard shows income, committed spending, remaining money, savings, goals, and joint account leftover.

Tasks:

- Add summary API.
- Add dashboard route.
- Add loading and error states.

### Manage member budgets

Acceptance criteria:

- User can view and edit recurring bills and savings for Neil and Katie.
- Member totals update after edits.

Tasks:

- Add member budget API.
- Add budget item update API.
- Add member budget screens.
- Add editable amount rows.

### Manage joint account

Acceptance criteria:

- User can view and edit shared bills.
- App shows total monthly cost, each-person share, contributions, and leftover.

Tasks:

- Add joint account API.
- Add joint bill update API.
- Add joint account screen.

### Track finance goals

Acceptance criteria:

- User can view and edit runway and emergency fund targets.
- Goal progress is shown visually.

Tasks:

- Add goals API.
- Add goal update API.
- Add goals screen.

### Add test and quality baseline

Acceptance criteria:

- Backend calculation tests pass.
- Frontend smoke test is available.
- README explains test commands and known local tooling requirements.

Tasks:

- Add Go tests for budget calculations.
- Add frontend smoke test.
- Add lint/test Makefile targets.

