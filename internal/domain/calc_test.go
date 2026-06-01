package domain

import "testing"

func TestBuildMemberBudgetMatchesNeilSpreadsheetTotals(t *testing.T) {
	member := Member{ID: 1, Name: "Neil"}
	items := []BudgetItem{
		{Kind: "bill", Label: "Joint Account", Amount: Pounds(1100)},
		{Kind: "bill", Label: "Car", Amount: Pounds(336.11)},
		{Kind: "bill", Label: "Credit Card", Amount: Pounds(0)},
		{Kind: "bill", Label: "David Lloyd", Amount: Pounds(189)},
		{Kind: "bill", Label: "NUFC", Amount: Pounds(63.62)},
		{Kind: "bill", Label: "Car Insurance", Amount: Pounds(21.10)},
		{Kind: "bill", Label: "Car Tax", Amount: Pounds(17.06)},
		{Kind: "bill", Label: "Phone", Amount: Pounds(29.96)},
		{Kind: "bill", Label: "Spotify", Amount: Pounds(14.99)},
		{Kind: "bill", Label: "O2", Amount: Pounds(9.79)},
		{Kind: "bill", Label: "Strava", Amount: Pounds(8.99)},
		{Kind: "bill", Label: "Amazon Prime", Amount: Pounds(8.99)},
		{Kind: "bill", Label: "BattleNet", Amount: Pounds(8.99)},
		{Kind: "saving", Label: "Emergency Fund", Amount: Pounds(300)},
		{Kind: "saving", Label: "Vanguard ISA", Amount: Pounds(500)},
		{Kind: "saving", Label: "Vanguard SIPP", Amount: Pounds(200)},
		{Kind: "saving", Label: "Joint Savings", Amount: Pounds(100)},
	}
	allocations := []AllocationRule{{Label: "70% - Living Expenses", Amount: Pounds(2516.42)}}

	got := BuildMemberBudget(member, Pounds(3594.88), items, allocations)
	if got.BillsTotal != Pounds(1808.60) {
		t.Fatalf("bills total = %s, want £1808.60", got.BillsTotal)
	}
	if got.CommittedTotal != Pounds(2908.60) {
		t.Fatalf("committed total = %s, want £2908.60", got.CommittedTotal)
	}
	if got.Remaining != Pounds(686.28) {
		t.Fatalf("remaining = %s, want £686.28", got.Remaining)
	}
}

func TestBuildJointAccountMatchesSpreadsheetTotals(t *testing.T) {
	joint := BuildJointAccount(nil, []JointAccountItem{
		{Label: "Mortgage", Amount: Pounds(1199.27)},
		{Label: "Council Tax", Amount: Pounds(179)},
	}, []JointContribution{
		{Member: "Neil", Amount: Pounds(1000)},
		{Member: "Katie", Amount: Pounds(850)},
	})

	if joint.Total != Pounds(1378.27) {
		t.Fatalf("joint total = %s, want £1378.27", joint.Total)
	}
	if joint.PerPerson != Pounds(689.14) {
		t.Fatalf("per person = %s, want £689.14", joint.PerPerson)
	}
	if joint.Leftover != Pounds(471.73) {
		t.Fatalf("leftover = %s, want £471.73", joint.Leftover)
	}
}
