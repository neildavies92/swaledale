package store

import (
	"context"
	"errors"
	"time"

	"github.com/neildavies/swaledale/internal/domain"
)

var ErrEmailTaken = errors.New("email already registered")

type Store interface {
	Close()
	RegisterUser(ctx context.Context, input RegisterUserInput) (domain.User, error)
	UserByEmail(ctx context.Context, email string) (domain.User, string, error)
	CreateSession(ctx context.Context, tokenHash string, userID int64, expiresAt time.Time) error
	UserBySession(ctx context.Context, tokenHash string) (domain.User, error)
	DeleteSession(ctx context.Context, tokenHash string) error
	Summary(ctx context.Context, householdID int64) (domain.Summary, error)
	Members(ctx context.Context, householdID int64) ([]domain.Member, error)
	MemberBudget(ctx context.Context, householdID int64, memberID int64) (domain.MemberBudget, error)
	UpdateBudgetItem(ctx context.Context, householdID int64, memberID int64, itemID int64, input UpdateBudgetItemInput) (domain.MemberBudget, error)
	JointAccount(ctx context.Context, householdID int64) (domain.JointAccount, error)
	UpdateJointAccountItem(ctx context.Context, householdID int64, itemID int64, input UpdateMoneyLabelInput) (domain.JointAccount, error)
	Goals(ctx context.Context, householdID int64) ([]domain.Goal, error)
	UpdateGoal(ctx context.Context, householdID int64, goalID int64, input UpdateGoalInput) ([]domain.Goal, error)
}

type RegisterUserInput struct {
	Name          string
	Email         string
	PasswordHash  string
	HouseholdName string
}

type UpdateMoneyLabelInput struct {
	Label  string       `json:"label"`
	Amount domain.Money `json:"amount"`
}

type UpdateBudgetItemInput = UpdateMoneyLabelInput

type UpdateGoalInput struct {
	Name    string       `json:"name"`
	Target  domain.Money `json:"target"`
	Current domain.Money `json:"current"`
	Notes   string       `json:"notes"`
}
