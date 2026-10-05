// Package legacybudget supports the existing budget API while FIRE replaces it.
// It is not the canonical financial model.
package legacybudget

import (
	"time"

	"github.com/neildavies/swaledale/internal/domain"
)

type Snapshot struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Month     time.Time `json:"month"`
	SourceURL string    `json:"sourceUrl"`
}

type IncomeEntry struct {
	ID       int64        `json:"id"`
	MemberID int64        `json:"memberId"`
	Label    string       `json:"label"`
	Amount   domain.Money `json:"amount"`
}

type BudgetItem struct {
	ID         int64        `json:"id"`
	MemberID   int64        `json:"memberId"`
	CategoryID int64        `json:"categoryId"`
	Category   string       `json:"category"`
	Kind       string       `json:"kind"`
	Label      string       `json:"label"`
	Amount     domain.Money `json:"amount"`
}

type AllocationRule struct {
	ID       int64        `json:"id"`
	MemberID int64        `json:"memberId"`
	Label    string       `json:"label"`
	Percent  int          `json:"percent"`
	Amount   domain.Money `json:"amount"`
}

type Goal struct {
	ID      int64        `json:"id"`
	Name    string       `json:"name"`
	Owner   string       `json:"owner"`
	Target  domain.Money `json:"target"`
	Current domain.Money `json:"current"`
	Notes   string       `json:"notes"`
}

type JointAccountItem struct {
	ID        int64        `json:"id"`
	Label     string       `json:"label"`
	Amount    domain.Money `json:"amount"`
	PerPerson domain.Money `json:"perPerson"`
}

type JointContribution struct {
	MemberID int64        `json:"memberId"`
	Member   string       `json:"member"`
	Amount   domain.Money `json:"amount"`
}

type MemberBudget struct {
	Member          domain.Member    `json:"member"`
	Income          domain.Money     `json:"income"`
	Items           []BudgetItem     `json:"items"`
	Allocations     []AllocationRule `json:"allocations"`
	BillsTotal      domain.Money     `json:"billsTotal"`
	SavingsTotal    domain.Money     `json:"savingsTotal"`
	CommittedTotal  domain.Money     `json:"committedTotal"`
	Remaining       domain.Money     `json:"remaining"`
	SpendableIncome domain.Money     `json:"spendableIncome"`
}

type JointAccount struct {
	Wages         []IncomeEntry       `json:"wages"`
	Items         []JointAccountItem  `json:"items"`
	Contributions []JointContribution `json:"contributions"`
	Total         domain.Money        `json:"total"`
	PerPerson     domain.Money        `json:"perPerson"`
	Leftover      domain.Money        `json:"leftover"`
}

type Summary struct {
	Household           domain.Household `json:"household"`
	Snapshot            Snapshot         `json:"snapshot"`
	Members             []MemberBudget   `json:"members"`
	Goals               []Goal           `json:"goals"`
	JointAccount        JointAccount     `json:"jointAccount"`
	TotalPersonalIncome domain.Money     `json:"totalPersonalIncome"`
	TotalCommitted      domain.Money     `json:"totalCommitted"`
	TotalRemaining      domain.Money     `json:"totalRemaining"`
	TotalSavings        domain.Money     `json:"totalSavings"`
}
