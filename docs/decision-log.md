# Decision Log

## Go API + React Frontend

Decision: Use Go for the API and React for the frontend.

Reason: The project is explicitly a Go learning project, while React keeps the UI flexible for dashboards and inline edits.

## Postgres As Source Of Truth

Decision: Store POC data in Postgres instead of reading directly from Google Sheets at runtime.

Reason: The app should become a proper application with structured data, migrations, and testable calculations.

## Dev-Only POC

Decision: Do not solve deployment or authentication yet.

Reason: The first milestone is proving spreadsheet parity and the household finance model.

## Fixed Spreadsheet Snapshot

Decision: Seed from the current spreadsheet values once.

Reason: Live sync adds complexity before the core model is useful.

## Future Considerations

- Add authentication before any cloud deployment.
- Add audit history for budget changes.
- Decide whether to support Google Sheets import/export.
- Explore Open Banking only after recurring budget workflows are solid.

