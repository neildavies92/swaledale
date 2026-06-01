-- +goose Up
CREATE TABLE households (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'GBP',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE members (
    id BIGSERIAL PRIMARY KEY,
    household_id BIGINT NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (household_id, name)
);

CREATE TABLE monthly_snapshots (
    id BIGSERIAL PRIMARY KEY,
    household_id BIGINT NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    month DATE NOT NULL,
    source_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE income_entries (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    member_id BIGINT REFERENCES members(id) ON DELETE CASCADE,
    context TEXT NOT NULL CHECK (context IN ('personal', 'joint')),
    label TEXT NOT NULL,
    amount_pence BIGINT NOT NULL
);

CREATE TABLE budget_categories (
    id BIGSERIAL PRIMARY KEY,
    household_id BIGINT NOT NULL REFERENCES households(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    kind TEXT NOT NULL CHECK (kind IN ('bill', 'saving', 'allocation')),
    UNIQUE (household_id, name, kind)
);

CREATE TABLE budget_items (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES budget_categories(id),
    label TEXT NOT NULL,
    amount_pence BIGINT NOT NULL
);

CREATE TABLE allocation_rules (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    percent INTEGER NOT NULL,
    amount_pence BIGINT NOT NULL
);

CREATE TABLE household_goals (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    owner_member_id BIGINT REFERENCES members(id) ON DELETE SET NULL,
    name TEXT NOT NULL,
    target_pence BIGINT NOT NULL,
    current_pence BIGINT NOT NULL DEFAULT 0,
    notes TEXT NOT NULL DEFAULT ''
);

CREATE TABLE joint_account_items (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    label TEXT NOT NULL,
    amount_pence BIGINT NOT NULL
);

CREATE TABLE joint_account_contributions (
    id BIGSERIAL PRIMARY KEY,
    snapshot_id BIGINT NOT NULL REFERENCES monthly_snapshots(id) ON DELETE CASCADE,
    member_id BIGINT NOT NULL REFERENCES members(id) ON DELETE CASCADE,
    amount_pence BIGINT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS joint_account_contributions;
DROP TABLE IF EXISTS joint_account_items;
DROP TABLE IF EXISTS household_goals;
DROP TABLE IF EXISTS allocation_rules;
DROP TABLE IF EXISTS budget_items;
DROP TABLE IF EXISTS budget_categories;
DROP TABLE IF EXISTS income_entries;
DROP TABLE IF EXISTS monthly_snapshots;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS households;

