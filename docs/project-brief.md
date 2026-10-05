# Swaledale project brief

Build a live, provider-independent household FIRE planning platform using Go, PostgreSQL, Grafana and automated financial ingestion. Start with one household; model any number of members and explicit personal/joint ownership without SaaS tenant/team abstractions.

The [FIRE Dashboard Project tab](https://docs.google.com/spreadsheets/d/1pStlDjjlz6qp1908SJvJVThnQWj3SH0DzKO70mDza4g/edit?gid=944913493#gid=944913493) owns the backlog, account inventory, architecture principle and acceptance criteria. Repository documentation describes implementation; it is not a competing plan.

Stable concepts are households, members, ownership, financial roles, account types, access classes and financial observations. Providers, products, investment instruments held by an account, and connector mechanisms can change independently of account identity.

FIRE-001 delivers canonical Go types, explicit validation and representative tests. It reuses the existing app, auth, PostgreSQL tooling and React foundation. Existing budget endpoints remain usable through an isolated compatibility package; spreadsheet parity is no longer a requirement. Spendable income now means income less committed bills and savings for every member. Joint budget per-person figures use the actual household member count and remain an illustrative equal split, not canonical account ownership.

Later tickets own persistence, provider configuration, baseline ingestion, household/member FIRE calculations and Grafana. React can become the management interface for members, accounts, connectors and assumptions. The fixed development seed is legacy test data, not the financial baseline for FIRE.

FIRE-001 excludes connector implementations, database migrations, forecasting, tax/pension access rules, transaction categorisation, product optimisation, market data, AWS/Terraform and frontend redesign. Existing authentication and household-scoped application queries are preserved.
