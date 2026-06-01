# Swaledale Frontend

React + TypeScript frontend for the Swaledale household finance POC.

## Stack

- React 19
- React Router 7
- TanStack Query
- Vite 8
- Tailwind CSS 4
- Vitest
- ESLint flat config

## Setup

From the repository root:

```sh
nvm use
cd frontend
npm install
npm run dev
```

The dev server runs on `http://localhost:5173` and proxies API requests to `http://localhost:8080`.

## Commands

```sh
npm run dev
npm test
npm run lint
npm run build
```

## Routes

- `/`: household dashboard
- `/members/:memberId`: member budget
- `/joint-account`: shared household bills
- `/goals`: runway and emergency targets

