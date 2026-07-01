package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neildavies/swaledale/internal/domain"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgres(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) RegisterUser(ctx context.Context, input RegisterUserInput) (domain.User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	var householdID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO households (name, currency)
		VALUES ($1, 'GBP')
		RETURNING id
	`, input.HouseholdName).Scan(&householdID); err != nil {
		return domain.User{}, err
	}

	var memberID int64
	if err := tx.QueryRow(ctx, `
		INSERT INTO members (household_id, name)
		VALUES ($1, $2)
		RETURNING id
	`, householdID, input.Name).Scan(&memberID); err != nil {
		return domain.User{}, err
	}

	month := time.Now().UTC()
	month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)
	if _, err := tx.Exec(ctx, `
		INSERT INTO monthly_snapshots (household_id, name, month, source_url)
		VALUES ($1, 'Starter Budget', $2, '')
	`, householdID, month); err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		HouseholdID: householdID,
		MemberID:    &memberID,
		Name:        input.Name,
		Email:       input.Email,
	}
	err = tx.QueryRow(ctx, `
		INSERT INTO users (household_id, member_id, name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, householdID, memberID, input.Name, input.Email, input.PasswordHash).Scan(&user.ID)
	if isUniqueViolation(err) {
		return domain.User{}, ErrEmailTaken
	}
	if err != nil {
		return domain.User{}, err
	}
	return user, tx.Commit(ctx)
}

func (s *PostgresStore) UserByEmail(ctx context.Context, email string) (domain.User, string, error) {
	var user domain.User
	var passwordHash string
	err := s.pool.QueryRow(ctx, `
		SELECT id, household_id, member_id, name, email, password_hash
		FROM users
		WHERE lower(email) = lower($1)
	`, email).Scan(&user.ID, &user.HouseholdID, &user.MemberID, &user.Name, &user.Email, &passwordHash)
	return user, passwordHash, err
}

func (s *PostgresStore) CreateSession(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		INSERT INTO sessions (token_hash, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, tokenHash, userID, expiresAt)
	return err
}

func (s *PostgresStore) UserBySession(ctx context.Context, tokenHash string) (domain.User, error) {
	var user domain.User
	err := s.pool.QueryRow(ctx, `
		SELECT u.id, u.household_id, u.member_id, u.name, u.email
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > now()
	`, tokenHash).Scan(&user.ID, &user.HouseholdID, &user.MemberID, &user.Name, &user.Email)
	return user, err
}

func (s *PostgresStore) DeleteSession(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}

func (s *PostgresStore) Summary(ctx context.Context, householdID int64) (domain.Summary, error) {
	household, snapshot, err := s.currentContext(ctx, householdID)
	if err != nil {
		return domain.Summary{}, err
	}
	members, err := s.Members(ctx, householdID)
	if err != nil {
		return domain.Summary{}, err
	}
	budgets := make([]domain.MemberBudget, 0, len(members))
	for _, member := range members {
		budget, err := s.MemberBudget(ctx, householdID, member.ID)
		if err != nil {
			return domain.Summary{}, err
		}
		budgets = append(budgets, budget)
	}
	goals, err := s.Goals(ctx, householdID)
	if err != nil {
		return domain.Summary{}, err
	}
	joint, err := s.JointAccount(ctx, householdID)
	if err != nil {
		return domain.Summary{}, err
	}
	return domain.BuildSummary(household, snapshot, budgets, goals, joint), nil
}

func (s *PostgresStore) Members(ctx context.Context, householdID int64) ([]domain.Member, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, name FROM members WHERE household_id = $1 ORDER BY id`, householdID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.Member
	for rows.Next() {
		var member domain.Member
		if err := rows.Scan(&member.ID, &member.Name); err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, rows.Err()
}

func (s *PostgresStore) MemberBudget(ctx context.Context, householdID int64, memberID int64) (domain.MemberBudget, error) {
	member, err := s.member(ctx, householdID, memberID)
	if err != nil {
		return domain.MemberBudget{}, err
	}
	snapshotID, err := s.currentSnapshotID(ctx, householdID)
	if err != nil {
		return domain.MemberBudget{}, err
	}
	income, err := s.personalIncome(ctx, snapshotID, memberID)
	if err != nil {
		return domain.MemberBudget{}, err
	}
	items, err := s.budgetItems(ctx, snapshotID, memberID)
	if err != nil {
		return domain.MemberBudget{}, err
	}
	allocations, err := s.allocationRules(ctx, snapshotID, memberID)
	if err != nil {
		return domain.MemberBudget{}, err
	}
	return domain.BuildMemberBudget(member, income, items, allocations), nil
}

func (s *PostgresStore) UpdateBudgetItem(ctx context.Context, householdID int64, memberID int64, itemID int64, input UpdateBudgetItemInput) (domain.MemberBudget, error) {
	if _, err := s.member(ctx, householdID, memberID); err != nil {
		return domain.MemberBudget{}, err
	}
	tag, err := s.pool.Exec(ctx, `
		UPDATE budget_items
		SET label = $3, amount_pence = $4
		WHERE member_id = $1 AND id = $2
	`, memberID, itemID, input.Label, int64(input.Amount))
	if err != nil {
		return domain.MemberBudget{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.MemberBudget{}, pgx.ErrNoRows
	}
	return s.MemberBudget(ctx, householdID, memberID)
}

func (s *PostgresStore) JointAccount(ctx context.Context, householdID int64) (domain.JointAccount, error) {
	snapshotID, err := s.currentSnapshotID(ctx, householdID)
	if err != nil {
		return domain.JointAccount{}, err
	}
	wages, err := s.jointWages(ctx, snapshotID)
	if err != nil {
		return domain.JointAccount{}, err
	}
	items, err := s.jointItems(ctx, snapshotID)
	if err != nil {
		return domain.JointAccount{}, err
	}
	contributions, err := s.jointContributions(ctx, snapshotID)
	if err != nil {
		return domain.JointAccount{}, err
	}
	return domain.BuildJointAccount(wages, items, contributions), nil
}

func (s *PostgresStore) UpdateJointAccountItem(ctx context.Context, householdID int64, itemID int64, input UpdateMoneyLabelInput) (domain.JointAccount, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE joint_account_items i
		SET label = $2, amount_pence = $3
		FROM monthly_snapshots s
		WHERE i.id = $1 AND s.id = i.snapshot_id AND s.household_id = $4
	`, itemID, input.Label, int64(input.Amount), householdID)
	if err != nil {
		return domain.JointAccount{}, err
	}
	if tag.RowsAffected() == 0 {
		return domain.JointAccount{}, pgx.ErrNoRows
	}
	return s.JointAccount(ctx, householdID)
}

func (s *PostgresStore) Goals(ctx context.Context, householdID int64) ([]domain.Goal, error) {
	snapshotID, err := s.currentSnapshotID(ctx, householdID)
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, `
		SELECT g.id, g.name, COALESCE(m.name, 'Household'), g.target_pence, g.current_pence, g.notes
		FROM household_goals g
		LEFT JOIN members m ON m.id = g.owner_member_id
		WHERE g.snapshot_id = $1
		ORDER BY g.id
	`, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []domain.Goal
	for rows.Next() {
		var goal domain.Goal
		var target, current int64
		if err := rows.Scan(&goal.ID, &goal.Name, &goal.Owner, &target, &current, &goal.Notes); err != nil {
			return nil, err
		}
		goal.Target = domain.Money(target)
		goal.Current = domain.Money(current)
		goals = append(goals, goal)
	}
	return goals, rows.Err()
}

func (s *PostgresStore) UpdateGoal(ctx context.Context, householdID int64, goalID int64, input UpdateGoalInput) ([]domain.Goal, error) {
	tag, err := s.pool.Exec(ctx, `
		UPDATE household_goals g
		SET name = $2, target_pence = $3, current_pence = $4, notes = $5
		FROM monthly_snapshots s
		WHERE g.id = $1 AND s.id = g.snapshot_id AND s.household_id = $6
	`, goalID, input.Name, int64(input.Target), int64(input.Current), input.Notes, householdID)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}
	return s.Goals(ctx, householdID)
}

func (s *PostgresStore) currentContext(ctx context.Context, householdID int64) (domain.Household, domain.Snapshot, error) {
	var household domain.Household
	var snapshot domain.Snapshot
	var month time.Time
	err := s.pool.QueryRow(ctx, `
		SELECT h.id, h.name, h.currency, s.id, s.name, s.month, s.source_url
		FROM monthly_snapshots s
		JOIN households h ON h.id = s.household_id
		WHERE h.id = $1
		ORDER BY s.month DESC, s.id DESC
		LIMIT 1
	`, householdID).Scan(&household.ID, &household.Name, &household.Currency, &snapshot.ID, &snapshot.Name, &month, &snapshot.SourceURL)
	if err != nil {
		return household, snapshot, fmt.Errorf("load current snapshot: %w", err)
	}
	snapshot.Month = month
	return household, snapshot, nil
}

func (s *PostgresStore) currentSnapshotID(ctx context.Context, householdID int64) (int64, error) {
	var id int64
	err := s.pool.QueryRow(ctx, `
		SELECT id FROM monthly_snapshots
		WHERE household_id = $1
		ORDER BY month DESC, id DESC
		LIMIT 1
	`, householdID).Scan(&id)
	return id, err
}

func (s *PostgresStore) member(ctx context.Context, householdID int64, memberID int64) (domain.Member, error) {
	var member domain.Member
	err := s.pool.QueryRow(ctx, `
		SELECT id, name FROM members WHERE id = $1 AND household_id = $2
	`, memberID, householdID).Scan(&member.ID, &member.Name)
	return member, err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

func (s *PostgresStore) personalIncome(ctx context.Context, snapshotID int64, memberID int64) (domain.Money, error) {
	var amount int64
	err := s.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_pence), 0)
		FROM income_entries
		WHERE snapshot_id = $1 AND member_id = $2 AND context = 'personal'
	`, snapshotID, memberID).Scan(&amount)
	return domain.Money(amount), err
}

func (s *PostgresStore) budgetItems(ctx context.Context, snapshotID int64, memberID int64) ([]domain.BudgetItem, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT i.id, i.member_id, c.id, c.name, c.kind, i.label, i.amount_pence
		FROM budget_items i
		JOIN budget_categories c ON c.id = i.category_id
		WHERE i.snapshot_id = $1 AND i.member_id = $2
		ORDER BY c.kind, i.id
	`, snapshotID, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.BudgetItem
	for rows.Next() {
		var item domain.BudgetItem
		var amount int64
		if err := rows.Scan(&item.ID, &item.MemberID, &item.CategoryID, &item.Category, &item.Kind, &item.Label, &amount); err != nil {
			return nil, err
		}
		item.Amount = domain.Money(amount)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *PostgresStore) allocationRules(ctx context.Context, snapshotID int64, memberID int64) ([]domain.AllocationRule, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, member_id, label, percent, amount_pence
		FROM allocation_rules
		WHERE snapshot_id = $1 AND member_id = $2
		ORDER BY id
	`, snapshotID, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []domain.AllocationRule
	for rows.Next() {
		var rule domain.AllocationRule
		var amount int64
		if err := rows.Scan(&rule.ID, &rule.MemberID, &rule.Label, &rule.Percent, &amount); err != nil {
			return nil, err
		}
		rule.Amount = domain.Money(amount)
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}

func (s *PostgresStore) jointWages(ctx context.Context, snapshotID int64) ([]domain.IncomeEntry, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT e.id, e.member_id, e.label, e.amount_pence
		FROM income_entries e
		WHERE e.snapshot_id = $1 AND e.context = 'joint'
		ORDER BY e.id
	`, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wages []domain.IncomeEntry
	for rows.Next() {
		var wage domain.IncomeEntry
		var memberID *int64
		var amount int64
		if err := rows.Scan(&wage.ID, &memberID, &wage.Label, &amount); err != nil {
			return nil, err
		}
		if memberID != nil {
			wage.MemberID = *memberID
		}
		wage.Amount = domain.Money(amount)
		wages = append(wages, wage)
	}
	return wages, rows.Err()
}

func (s *PostgresStore) jointItems(ctx context.Context, snapshotID int64) ([]domain.JointAccountItem, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, label, amount_pence
		FROM joint_account_items
		WHERE snapshot_id = $1
		ORDER BY id
	`, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.JointAccountItem
	for rows.Next() {
		var item domain.JointAccountItem
		var amount int64
		if err := rows.Scan(&item.ID, &item.Label, &amount); err != nil {
			return nil, err
		}
		item.Amount = domain.Money(amount)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *PostgresStore) jointContributions(ctx context.Context, snapshotID int64) ([]domain.JointContribution, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT c.member_id, m.name, c.amount_pence
		FROM joint_account_contributions c
		JOIN members m ON m.id = c.member_id
		WHERE c.snapshot_id = $1
		ORDER BY c.id
	`, snapshotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributions []domain.JointContribution
	for rows.Next() {
		var contribution domain.JointContribution
		var amount int64
		if err := rows.Scan(&contribution.MemberID, &contribution.Member, &amount); err != nil {
			return nil, err
		}
		contribution.Amount = domain.Money(amount)
		contributions = append(contributions, contribution)
	}
	return contributions, rows.Err()
}
