# Decision log

## 2026-10-05 — FIRE-001 canonical household domain

Swaledale now targets household financial independence planning. The [FIRE Dashboard Project tab](https://docs.google.com/spreadsheets/d/1pStlDjjlz6qp1908SJvJVThnQWj3SH0DzKO70mDza4g/edit?gid=944913493#gid=944913493) supersedes the old Jira/spreadsheet-parity backlog, which has been removed.

Reuse Go/chi/pgx/PostgreSQL, goose/sqlc tooling, React/Vite, household-scoped authentication, Household, Member and integer Money. No tenant/team abstraction is added. Member now carries the household key already present in PostgreSQL; store reads populate it.

Account identity is independent of providers and instruments. ProviderConnection models effective-dated mappings; Holding models a time-stamped account/instrument relationship. FinancialRole is extensible configuration; AccountType and AccessClass are validated typed values. Pension eligibility ages remain calculation assumptions, not account fields.

Ownership uses complete sets of integer basis-point shares, unique owners and household validation. Persistence must version ownership sets before ownership edits are offered. Snapshots retain effective/recording timestamps, currency and provenance; persistence must preserve history and source references.

Money remains integer minor units with explicit currency in Amount. The canonical wire shape is integer minorUnits plus currency. Legacy decimal-GBP JSON is preserved using exact parsing/formatting; floating-point helpers and seed conversions are removed. Quantities use exact decimal strings. No FX engine is introduced.

The active budget API types/calculations move to internal/legacybudget. Allocation labels no longer drive calculations, every member's spendable income is income minus committed bills/savings, and per-person budget estimates use actual member count. This intentionally changes spreadsheet-specific spendable figures. No React rewrite or destructive migration is needed. Unused BudgetCategory is removed.

FIRE-003 owns migrations and historical storage constraints; FIRE-024 owns provider configuration; FIRE-005 owns calculations; FIRE-015 owns holdings batch completeness. FIRE-027 tracks eventual budget API/UI retirement and FIRE-028 tracks household scoping of unused sqlc templates before adoption.

## Retained foundations and historical decisions

Go remains the backend learning focus, React the flexible management UI, and PostgreSQL the structured financial data store. The initial fixed spreadsheet seed remains development-only reference data. Earlier decisions to defer authentication and target spreadsheet parity are superseded: authentication already exists and parity is no longer the objective.
