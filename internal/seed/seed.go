package seed

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neildavies/swaledale/internal/auth"
	"github.com/neildavies/swaledale/internal/domain"
)

const SourceURL = "https://docs.google.com/spreadsheets/d/1pStlDjjlz6qp1908SJvJVThnQWj3SH0DzKO70mDza4g/edit"

// DevPassword is the shared password for the seeded dev accounts.
const DevPassword = "swaledale-dev"

type memberRef struct {
	ID   int64
	Name string
}

func Run(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := reset(ctx, tx); err != nil {
		return err
	}

	householdID, err := insertHousehold(ctx, tx)
	if err != nil {
		return err
	}
	members, err := insertMembers(ctx, tx, householdID)
	if err != nil {
		return err
	}
	snapshotID, err := insertSnapshot(ctx, tx, householdID)
	if err != nil {
		return err
	}
	categories, err := insertCategories(ctx, tx, householdID)
	if err != nil {
		return err
	}
	if err := insertIncome(ctx, tx, snapshotID, members); err != nil {
		return err
	}
	if err := insertBudgetItems(ctx, tx, snapshotID, members, categories); err != nil {
		return err
	}
	if err := insertAllocations(ctx, tx, snapshotID, members); err != nil {
		return err
	}
	if err := insertGoals(ctx, tx, snapshotID, members); err != nil {
		return err
	}
	if err := insertJointAccount(ctx, tx, snapshotID, members); err != nil {
		return err
	}
	if err := insertUsers(ctx, tx, householdID, members); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func reset(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		TRUNCATE households RESTART IDENTITY CASCADE
	`)
	return err
}

func insertHousehold(ctx context.Context, tx pgx.Tx) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		INSERT INTO households (name, currency)
		VALUES ('Swaledale Household', 'GBP')
		RETURNING id
	`).Scan(&id)
	return id, err
}

func insertMembers(ctx context.Context, tx pgx.Tx, householdID int64) (map[string]memberRef, error) {
	members := map[string]memberRef{}
	for _, name := range []string{"Neil", "Katie"} {
		var id int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO members (household_id, name)
			VALUES ($1, $2)
			RETURNING id
		`, householdID, name).Scan(&id); err != nil {
			return nil, err
		}
		members[name] = memberRef{ID: id, Name: name}
	}
	return members, nil
}

func insertSnapshot(ctx context.Context, tx pgx.Tx, householdID int64) (int64, error) {
	var id int64
	month := time.Date(2026, time.May, 1, 0, 0, 0, 0, time.UTC)
	err := tx.QueryRow(ctx, `
		INSERT INTO monthly_snapshots (household_id, name, month, source_url)
		VALUES ($1, 'Spreadsheet POC Snapshot', $2, $3)
		RETURNING id
	`, householdID, month, SourceURL).Scan(&id)
	return id, err
}

func insertCategories(ctx context.Context, tx pgx.Tx, householdID int64) (map[string]int64, error) {
	type category struct {
		name string
		kind string
	}
	categories := []category{
		{name: "Bills", kind: "bill"},
		{name: "Savings / Investments", kind: "saving"},
	}
	ids := map[string]int64{}
	for _, category := range categories {
		var id int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO budget_categories (household_id, name, kind)
			VALUES ($1, $2, $3)
			RETURNING id
		`, householdID, category.name, category.kind).Scan(&id); err != nil {
			return nil, err
		}
		ids[category.kind] = id
	}
	return ids, nil
}

func insertIncome(ctx context.Context, tx pgx.Tx, snapshotID int64, members map[string]memberRef) error {
	rows := []struct {
		member  string
		context string
		label   string
		amount  domain.Money
	}{
		{member: "Neil", context: "personal", label: "Income", amount: 359488},
		{member: "Katie", context: "personal", label: "Bills", amount: 191316},
		{member: "Neil", context: "joint", label: "Neil", amount: 361562},
		{member: "Katie", context: "joint", label: "Katie", amount: 191336},
	}
	for _, row := range rows {
		member := members[row.member]
		if _, err := tx.Exec(ctx, `
			INSERT INTO income_entries (snapshot_id, member_id, context, label, amount_pence)
			VALUES ($1, $2, $3, $4, $5)
		`, snapshotID, member.ID, row.context, row.label, int64(row.amount)); err != nil {
			return err
		}
	}
	return nil
}

func insertBudgetItems(ctx context.Context, tx pgx.Tx, snapshotID int64, members map[string]memberRef, categories map[string]int64) error {
	rows := []struct {
		member string
		kind   string
		label  string
		amount domain.Money
	}{
		{member: "Neil", kind: "bill", label: "Joint Account", amount: 110000},
		{member: "Neil", kind: "bill", label: "Car", amount: 33611},
		{member: "Neil", kind: "bill", label: "Credit Card", amount: 0},
		{member: "Neil", kind: "bill", label: "David Lloyd", amount: 18900},
		{member: "Neil", kind: "bill", label: "NUFC", amount: 6362},
		{member: "Neil", kind: "bill", label: "Car Insurance", amount: 2110},
		{member: "Neil", kind: "bill", label: "Car Tax", amount: 1706},
		{member: "Neil", kind: "bill", label: "Phone", amount: 2996},
		{member: "Neil", kind: "bill", label: "Spotify", amount: 1499},
		{member: "Neil", kind: "bill", label: "O2", amount: 979},
		{member: "Neil", kind: "bill", label: "Strava", amount: 899},
		{member: "Neil", kind: "bill", label: "Amazon Prime", amount: 899},
		{member: "Neil", kind: "bill", label: "BattleNet", amount: 899},
		{member: "Neil", kind: "saving", label: "Emergency Fund", amount: 30000},
		{member: "Neil", kind: "saving", label: "Vanguard ISA", amount: 50000},
		{member: "Neil", kind: "saving", label: "Vanguard SIPP", amount: 20000},
		{member: "Neil", kind: "saving", label: "Joint Savings", amount: 10000},
		{member: "Neil", kind: "saving", label: "Personal Savings", amount: 0},
		{member: "Katie", kind: "bill", label: "Joint Savings", amount: 0},
		{member: "Katie", kind: "saving", label: "Savings", amount: 30000},
		{member: "Katie", kind: "bill", label: "Car Tax", amount: 1443},
		{member: "Katie", kind: "bill", label: "Gym", amount: 2999},
		{member: "Katie", kind: "bill", label: "Car Insurance", amount: 2800},
		{member: "Katie", kind: "bill", label: "Joint Account", amount: 85000},
		{member: "Katie", kind: "bill", label: "Phone", amount: 979},
		{member: "Katie", kind: "bill", label: "Spotify", amount: 2199},
		{member: "Katie", kind: "bill", label: "Credit card", amount: 4794},
		{member: "Katie", kind: "bill", label: "Uber One", amount: 499},
		{member: "Katie", kind: "bill", label: "Dentist loan", amount: 5320},
		{member: "Katie", kind: "bill", label: "Petrol", amount: 5000},
	}
	for _, row := range rows {
		member := members[row.member]
		categoryID, ok := categories[row.kind]
		if !ok {
			return fmt.Errorf("missing category %q", row.kind)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO budget_items (snapshot_id, member_id, category_id, label, amount_pence)
			VALUES ($1, $2, $3, $4, $5)
		`, snapshotID, member.ID, categoryID, row.label, int64(row.amount)); err != nil {
			return err
		}
	}
	return nil
}

func insertAllocations(ctx context.Context, tx pgx.Tx, snapshotID int64, members map[string]memberRef) error {
	rows := []struct {
		member  string
		label   string
		percent int
		amount  domain.Money
	}{
		{member: "Neil", label: "10% - PYF", percent: 10, amount: 35949},
		{member: "Neil", label: "20% - Debts / Investment", percent: 20, amount: 71898},
		{member: "Neil", label: "70% - Living Expenses", percent: 70, amount: 251642},
		{member: "Neil", label: "Spendable Income", percent: 0, amount: 70782},
		{member: "Katie", label: "Bills/Expenses - 50%", percent: 50, amount: 95658},
		{member: "Katie", label: "Spendable Income - 30%", percent: 30, amount: 57395},
		{member: "Katie", label: "Savings/Investment - 20%", percent: 20, amount: 38263},
	}
	for _, row := range rows {
		member := members[row.member]
		if _, err := tx.Exec(ctx, `
			INSERT INTO allocation_rules (snapshot_id, member_id, label, percent, amount_pence)
			VALUES ($1, $2, $3, $4, $5)
		`, snapshotID, member.ID, row.label, row.percent, int64(row.amount)); err != nil {
			return err
		}
	}
	return nil
}

func insertGoals(ctx context.Context, tx pgx.Tx, snapshotID int64, members map[string]memberRef) error {
	rows := []struct {
		owner   string
		name    string
		target  domain.Money
		current domain.Money
		notes   string
	}{
		{owner: "Neil", name: "Runway Goal (Bills x 3 Months)", target: 542580, notes: "Neil bills x 3 months"},
		{owner: "Neil", name: "Emergency Fund", target: 500000, notes: "Cash emergency fund target"},
		{owner: "Neil", name: "Runway + Emergency", target: 1042580, notes: "Runway plus emergency target"},
		{owner: "Katie", name: "Runway Goal (Bills x 3 Months)", target: 333099, notes: "Katie bills x 3 months"},
		{owner: "Katie", name: "Emergency Fund", target: 500000, notes: "Cash emergency fund target"},
		{owner: "Katie", name: "Runway + Emergency", target: 833099, notes: "Runway plus emergency target"},
	}
	for _, row := range rows {
		member := members[row.owner]
		if _, err := tx.Exec(ctx, `
			INSERT INTO household_goals (snapshot_id, owner_member_id, name, target_pence, current_pence, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, snapshotID, member.ID, row.name, int64(row.target), int64(row.current), row.notes); err != nil {
			return err
		}
	}
	return nil
}

func insertUsers(ctx context.Context, tx pgx.Tx, householdID int64, members map[string]memberRef) error {
	passwordHash, err := auth.HashPassword(DevPassword)
	if err != nil {
		return err
	}
	rows := []struct {
		member string
		email  string
	}{
		{member: "Neil", email: "neil@swaledale.local"},
		{member: "Katie", email: "katie@swaledale.local"},
	}
	for _, row := range rows {
		member := members[row.member]
		if _, err := tx.Exec(ctx, `
			INSERT INTO users (household_id, member_id, name, email, password_hash)
			VALUES ($1, $2, $3, $4, $5)
		`, householdID, member.ID, member.Name, row.email, passwordHash); err != nil {
			return err
		}
	}
	return nil
}

func insertJointAccount(ctx context.Context, tx pgx.Tx, snapshotID int64, members map[string]memberRef) error {
	items := []struct {
		label  string
		amount domain.Money
	}{
		{label: "Mortgage", amount: 119927},
		{label: "Council Tax", amount: 17900},
		{label: "Shopping", amount: 15000},
		{label: "Gas / Electric", amount: 10042},
		{label: "Water", amount: 6041},
		{label: "Sofa", amount: 5175},
		{label: "Internet", amount: 3506},
		{label: "House Insurance", amount: 1875},
		{label: "Life Cover", amount: 1789},
		{label: "TV", amount: 1503},
		{label: "Window Cleaner", amount: 1600},
		{label: "Pet Insurance", amount: 881},
	}
	for _, item := range items {
		if _, err := tx.Exec(ctx, `
			INSERT INTO joint_account_items (snapshot_id, label, amount_pence)
			VALUES ($1, $2, $3)
		`, snapshotID, item.label, int64(item.amount)); err != nil {
			return err
		}
	}

	contributions := []struct {
		member string
		amount domain.Money
	}{
		{member: "Neil", amount: 100000},
		{member: "Katie", amount: 85000},
	}
	for _, contribution := range contributions {
		member := members[contribution.member]
		if _, err := tx.Exec(ctx, `
			INSERT INTO joint_account_contributions (snapshot_id, member_id, amount_pence)
			VALUES ($1, $2, $3)
		`, snapshotID, member.ID, int64(contribution.amount)); err != nil {
			return err
		}
	}
	return nil
}
