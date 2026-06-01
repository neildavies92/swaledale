DATABASE_URL ?= postgres://swaledale:swaledale@localhost:5432/swaledale?sslmode=disable
GOCACHE ?= $(CURDIR)/.cache/go-build
GOMODCACHE ?= $(CURDIR)/.cache/go-mod
NPM_CONFIG_CACHE ?= $(CURDIR)/.cache/npm
NODE_BIN ?= $(HOME)/.nvm/versions/node/v24.16.0/bin
GO_ENV = GOCACHE="$(GOCACHE)" GOMODCACHE="$(GOMODCACHE)"
NPM_ENV = PATH="$(NODE_BIN):$(PATH)" NPM_CONFIG_CACHE="$(NPM_CONFIG_CACHE)"
GOOSE = $(GO_ENV) go run github.com/pressly/goose/v3/cmd/goose@v3.24.1

.PHONY: dev api frontend db-up db-down migrate seed test lint

dev:
	$(MAKE) db-up
	@printf "Run the API and frontend in separate shells:\n  make api\n  make frontend\n"

api:
	$(GO_ENV) DATABASE_URL="$(DATABASE_URL)" go run ./cmd/api

frontend:
	cd frontend && $(NPM_ENV) npm run dev

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate:
	$(GOOSE) -dir migrations postgres "$(DATABASE_URL)" up

seed:
	$(GO_ENV) DATABASE_URL="$(DATABASE_URL)" go run ./cmd/seed

test:
	$(GO_ENV) go test ./cmd/... ./internal/...
	cd frontend && $(NPM_ENV) npm test

lint:
	$(GO_ENV) go test ./cmd/... ./internal/...
	cd frontend && $(NPM_ENV) npm run lint
