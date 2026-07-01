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
		{member: "Neil", context: "personal", label: "Income", amount: domain.Pounds(3594.88)},
		{member: "Katie", context: "personal", label: "Bills", amount: domain.Pounds(1913.16)},
		{member: "Neil", context: "joint", label: "Neil", amount: domain.Pounds(3615.62)},
		{member: "Katie", context: "joint", label: "Katie", amount: domain.Pounds(1913.36)},
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
		{member: "Neil", kind: "bill", label: "Joint Account", amount: domain.Pounds(1100)},
		{member: "Neil", kind: "bill", label: "Car", amount: domain.Pounds(336.11)},
		{member: "Neil", kind: "bill", label: "Credit Card", amount: domain.Pounds(0)},
		{member: "Neil", kind: "bill", label: "David Lloyd", amount: domain.Pounds(189)},
		{member: "Neil", kind: "bill", label: "NUFC", amount: domain.Pounds(63.62)},
		{member: "Neil", kind: "bill", label: "Car Insurance", amount: domain.Pounds(21.10)},
		{member: "Neil", kind: "bill", label: "Car Tax", amount: domain.Pounds(17.06)},
		{member: "Neil", kind: "bill", label: "Phone", amount: domain.Pounds(29.96)},
		{member: "Neil", kind: "bill", label: "Spotify", amount: domain.Pounds(14.99)},
		{member: "Neil", kind: "bill", label: "O2", amount: domain.Pounds(9.79)},
		{member: "Neil", kind: "bill", label: "Strava", amount: domain.Pounds(8.99)},
		{member: "Neil", kind: "bill", label: "Amazon Prime", amount: domain.Pounds(8.99)},
		{member: "Neil", kind: "bill", label: "BattleNet", amount: domain.Pounds(8.99)},
		{member: "Neil", kind: "saving", label: "Emergency Fund", amount: domain.Pounds(300)},
		{member: "Neil", kind: "saving", label: "Vanguard ISA", amount: domain.Pounds(500)},
		{member: "Neil", kind: "saving", label: "Vanguard SIPP", amount: domain.Pounds(200)},
		{member: "Neil", kind: "saving", label: "Joint Savings", amount: domain.Pounds(100)},
		{member: "Neil", kind: "saving", label: "Personal Savings", amount: domain.Pounds(0)},
		{member: "Katie", kind: "bill", label: "Joint Savings", amount: domain.Pounds(0)},
		{member: "Katie", kind: "saving", label: "Savings", amount: domain.Pounds(300)},
		{member: "Katie", kind: "bill", label: "Car Tax", amount: domain.Pounds(14.43)},
		{member: "Katie", kind: "bill", label: "Gym", amount: domain.Pounds(29.99)},
		{member: "Katie", kind: "bill", label: "Car Insurance", amount: domain.Pounds(28)},
		{member: "Katie", kind: "bill", label: "Joint Account", amount: domain.Pounds(850)},
		{member: "Katie", kind: "bill", label: "Phone", amount: domain.Pounds(9.79)},
		{member: "Katie", kind: "bill", label: "Spotify", amount: domain.Pounds(21.99)},
		{member: "Katie", kind: "bill", label: "Credit card", amount: domain.Pounds(47.94)},
		{member: "Katie", kind: "bill", label: "Uber One", amount: domain.Pounds(4.99)},
		{member: "Katie", kind: "bill", label: "Dentist loan", amount: domain.Pounds(53.20)},
		{member: "Katie", kind: "bill", label: "Petrol", amount: domain.Pounds(50)},
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
		{member: "Neil", label: "10% - PYF", percent: 10, amount: domain.Pounds(359.49)},
		{member: "Neil", label: "20% - Debts / Investment", percent: 20, amount: domain.Pounds(718.98)},
		{member: "Neil", label: "70% - Living Expenses", percent: 70, amount: domain.Pounds(2516.42)},
		{member: "Neil", label: "Spendable Income", percent: 0, amount: domain.Pounds(707.82)},
		{member: "Katie", label: "Bills/Expenses - 50%", percent: 50, amount: domain.Pounds(956.58)},
		{member: "Katie", label: "Spendable Income - 30%", percent: 30, amount: domain.Pounds(573.95)},
		{member: "Katie", label: "Savings/Investment - 20%", percent: 20, amount: domain.Pounds(382.63)},
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
		{owner: "Neil", name: "Runway Goal (Bills x 3 Months)", target: domain.Pounds(5425.80), notes: "Neil bills x 3 months"},
		{owner: "Neil", name: "Emergency Fund", target: domain.Pounds(5000), notes: "Cash emergency fund target"},
		{owner: "Neil", name: "Runway + Emergency", target: domain.Pounds(10425.80), notes: "Runway plus emergency target"},
		{owner: "Katie", name: "Runway Goal (Bills x 3 Months)", target: domain.Pounds(3330.99), notes: "Katie bills x 3 months"},
		{owner: "Katie", name: "Emergency Fund", target: domain.Pounds(5000), notes: "Cash emergency fund target"},
		{owner: "Katie", name: "Runway + Emergency", target: domain.Pounds(8330.99), notes: "Runway plus emergency target"},
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
		{label: "Mortgage", amount: domain.Pounds(1199.27)},
		{label: "Council Tax", amount: domain.Pounds(179)},
		{label: "Shopping", amount: domain.Pounds(150)},
		{label: "Gas / Electric", amount: domain.Pounds(100.42)},
		{label: "Water", amount: domain.Pounds(60.41)},
		{label: "Sofa", amount: domain.Pounds(51.75)},
		{label: "Internet", amount: domain.Pounds(35.06)},
		{label: "House Insurance", amount: domain.Pounds(18.75)},
		{label: "Life Cover", amount: domain.Pounds(17.89)},
		{label: "TV", amount: domain.Pounds(15.03)},
		{label: "Window Cleaner", amount: domain.Pounds(16)},
		{label: "Pet Insurance", amount: domain.Pounds(8.81)},
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
		{member: "Neil", amount: domain.Pounds(1000)},
		{member: "Katie", amount: domain.Pounds(850)},
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
