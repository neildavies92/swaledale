package legacybudget

import "github.com/neildavies/swaledale/internal/domain"

func SumBudgetItems(items []BudgetItem, kind string) domain.Money {
	var total domain.Money
	for _, item := range items {
		if kind == "" || item.Kind == kind {
			total += item.Amount
		}
	}
	return total
}

func BuildMemberBudget(member domain.Member, income domain.Money, items []BudgetItem, allocations []AllocationRule) MemberBudget {
	bills := SumBudgetItems(items, "bill")
	savings := SumBudgetItems(items, "saving")
	committed := bills + savings
	return MemberBudget{
		Member:          member,
		Income:          income,
		Items:           items,
		Allocations:     allocations,
		BillsTotal:      bills,
		SavingsTotal:    savings,
		CommittedTotal:  committed,
		Remaining:       income - committed,
		SpendableIncome: income - committed,
	}
}

func BuildJointAccount(wages []IncomeEntry, items []JointAccountItem, contributions []JointContribution, memberCount int) JointAccount {
	var total domain.Money
	for i := range items {
		items[i].PerPerson = divideRounded(items[i].Amount, domain.Money(memberCount))
		total += items[i].Amount
	}
	var contributionsTotal domain.Money
	for _, contribution := range contributions {
		contributionsTotal += contribution.Amount
	}
	return JointAccount{
		Wages:         wages,
		Items:         items,
		Contributions: contributions,
		Total:         total,
		PerPerson:     divideRounded(total, domain.Money(memberCount)),
		Leftover:      contributionsTotal - total,
	}
}

func BuildSummary(household domain.Household, snapshot Snapshot, members []MemberBudget, goals []Goal, joint JointAccount) Summary {
	var income, committed, remaining, savings domain.Money
	for _, member := range members {
		income += member.Income
		committed += member.CommittedTotal
		remaining += member.Remaining
		savings += member.SavingsTotal
	}
	return Summary{
		Household:           household,
		Snapshot:            snapshot,
		Members:             members,
		Goals:               goals,
		JointAccount:        joint,
		TotalPersonalIncome: income,
		TotalCommitted:      committed,
		TotalRemaining:      remaining,
		TotalSavings:        savings,
	}
}

func divideRounded(value domain.Money, divisor domain.Money) domain.Money {
	if divisor <= 0 {
		return 0
	}
	if value >= 0 {
		return (value + divisor/2) / divisor
	}
	return (value - divisor/2) / divisor
}
