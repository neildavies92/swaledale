package store

import (
	"context"

	"github.com/neildavies/swaledale/internal/domain"
)

type Store interface {
	Close()
	Summary(ctx context.Context) (domain.Summary, error)
	Members(ctx context.Context) ([]domain.Member, error)
	MemberBudget(ctx context.Context, memberID int64) (domain.MemberBudget, error)
	UpdateBudgetItem(ctx context.Context, memberID int64, itemID int64, input UpdateBudgetItemInput) (domain.MemberBudget, error)
	JointAccount(ctx context.Context) (domain.JointAccount, error)
	UpdateJointAccountItem(ctx context.Context, itemID int64, input UpdateMoneyLabelInput) (domain.JointAccount, error)
	Goals(ctx context.Context) ([]domain.Goal, error)
	UpdateGoal(ctx context.Context, goalID int64, input UpdateGoalInput) ([]domain.Goal, error)
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
