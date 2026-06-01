package domain

import "time"

type Household struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type Snapshot struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Month     time.Time `json:"month"`
	SourceURL string    `json:"sourceUrl"`
}

type Member struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type IncomeEntry struct {
	ID       int64  `json:"id"`
	MemberID int64  `json:"memberId"`
	Label    string `json:"label"`
	Amount   Money  `json:"amount"`
}

type BudgetCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type BudgetItem struct {
	ID         int64  `json:"id"`
	MemberID   int64  `json:"memberId"`
	CategoryID int64  `json:"categoryId"`
	Category   string `json:"category"`
	Kind       string `json:"kind"`
	Label      string `json:"label"`
	Amount     Money  `json:"amount"`
}

type AllocationRule struct {
	ID       int64  `json:"id"`
	MemberID int64  `json:"memberId"`
	Label    string `json:"label"`
	Percent  int    `json:"percent"`
	Amount   Money  `json:"amount"`
}

type Goal struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Target  Money  `json:"target"`
	Current Money  `json:"current"`
	Notes   string `json:"notes"`
}

type JointAccountItem struct {
	ID        int64  `json:"id"`
	Label     string `json:"label"`
	Amount    Money  `json:"amount"`
	PerPerson Money  `json:"perPerson"`
}

type JointContribution struct {
	MemberID int64  `json:"memberId"`
	Member   string `json:"member"`
	Amount   Money  `json:"amount"`
}

type MemberBudget struct {
	Member          Member           `json:"member"`
	Income          Money            `json:"income"`
	Items           []BudgetItem     `json:"items"`
	Allocations     []AllocationRule `json:"allocations"`
	BillsTotal      Money            `json:"billsTotal"`
	SavingsTotal    Money            `json:"savingsTotal"`
	CommittedTotal  Money            `json:"committedTotal"`
	Remaining       Money            `json:"remaining"`
	SpendableIncome Money            `json:"spendableIncome"`
}

type JointAccount struct {
	Wages         []IncomeEntry       `json:"wages"`
	Items         []JointAccountItem  `json:"items"`
	Contributions []JointContribution `json:"contributions"`
	Total         Money               `json:"total"`
	PerPerson     Money               `json:"perPerson"`
	Leftover      Money               `json:"leftover"`
}

type Summary struct {
	Household           Household      `json:"household"`
	Snapshot            Snapshot       `json:"snapshot"`
	Members             []MemberBudget `json:"members"`
	Goals               []Goal         `json:"goals"`
	JointAccount        JointAccount   `json:"jointAccount"`
	TotalPersonalIncome Money          `json:"totalPersonalIncome"`
	TotalCommitted      Money          `json:"totalCommitted"`
	TotalRemaining      Money          `json:"totalRemaining"`
	TotalSavings        Money          `json:"totalSavings"`
}
