package domain

func SumBudgetItems(items []BudgetItem, kind string) Money {
	var total Money
	for _, item := range items {
		if kind == "" || item.Kind == kind {
			total += item.Amount
		}
	}
	return total
}

func BuildMemberBudget(member Member, income Money, items []BudgetItem, allocations []AllocationRule) MemberBudget {
	bills := SumBudgetItems(items, "bill")
	savings := SumBudgetItems(items, "saving")
	committed := bills + savings
	spendable := Money(0)
	for _, allocation := range allocations {
		if allocation.Label == "Spendable Income" || allocation.Label == "Spendable Income - 30%" || allocation.Label == "70% - Living Expenses" {
			spendable = allocation.Amount
		}
	}
	if member.Name == "Neil" {
		living := allocationAmount(allocations, "70% - Living Expenses")
		if living != 0 {
			spendable = living - bills
		}
	}
	return MemberBudget{
		Member:          member,
		Income:          income,
		Items:           items,
		Allocations:     allocations,
		BillsTotal:      bills,
		SavingsTotal:    savings,
		CommittedTotal:  committed,
		Remaining:       income - committed,
		SpendableIncome: spendable,
	}
}

func BuildJointAccount(wages []IncomeEntry, items []JointAccountItem, contributions []JointContribution) JointAccount {
	var total Money
	for i := range items {
		items[i].PerPerson = divideRounded(items[i].Amount, 2)
		total += items[i].Amount
	}
	var contributionsTotal Money
	for _, contribution := range contributions {
		contributionsTotal += contribution.Amount
	}
	return JointAccount{
		Wages:         wages,
		Items:         items,
		Contributions: contributions,
		Total:         total,
		PerPerson:     divideRounded(total, 2),
		Leftover:      contributionsTotal - total,
	}
}

func BuildSummary(household Household, snapshot Snapshot, members []MemberBudget, goals []Goal, joint JointAccount) Summary {
	var income, committed, remaining, savings Money
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

func divideRounded(value Money, divisor Money) Money {
	if divisor == 0 {
		return 0
	}
	if value >= 0 {
		return (value + divisor/2) / divisor
	}
	return (value - divisor/2) / divisor
}

func allocationAmount(allocations []AllocationRule, label string) Money {
	for _, allocation := range allocations {
		if allocation.Label == label {
			return allocation.Amount
		}
	}
	return 0
}
