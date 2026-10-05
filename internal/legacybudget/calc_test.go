package legacybudget

import "github.com/neildavies/swaledale/internal/domain"

import "testing"

func TestBuildMemberBudgetMatchesNeilSpreadsheetTotals(t *testing.T) {
	member := domain.Member{ID: 1, Name: "Neil"}
	items := []BudgetItem{
		{Kind: "bill", Label: "Joint Account", Amount: 110000},
		{Kind: "bill", Label: "Car", Amount: 33611},
		{Kind: "bill", Label: "Credit Card", Amount: 0},
		{Kind: "bill", Label: "David Lloyd", Amount: 18900},
		{Kind: "bill", Label: "NUFC", Amount: 6362},
		{Kind: "bill", Label: "Car Insurance", Amount: 2110},
		{Kind: "bill", Label: "Car Tax", Amount: 1706},
		{Kind: "bill", Label: "Phone", Amount: 2996},
		{Kind: "bill", Label: "Spotify", Amount: 1499},
		{Kind: "bill", Label: "O2", Amount: 979},
		{Kind: "bill", Label: "Strava", Amount: 899},
		{Kind: "bill", Label: "Amazon Prime", Amount: 899},
		{Kind: "bill", Label: "BattleNet", Amount: 899},
		{Kind: "saving", Label: "Emergency Fund", Amount: 30000},
		{Kind: "saving", Label: "Vanguard ISA", Amount: 50000},
		{Kind: "saving", Label: "Vanguard SIPP", Amount: 20000},
		{Kind: "saving", Label: "Joint Savings", Amount: 10000},
	}
	allocations := []AllocationRule{{Label: "70% - Living Expenses", Amount: 251642}}

	got := BuildMemberBudget(member, 359488, items, allocations)
	if got.BillsTotal != 180860 {
		t.Fatalf("bills total = %s, want £1808.60", got.BillsTotal)
	}
	if got.CommittedTotal != 290860 {
		t.Fatalf("committed total = %s, want £2908.60", got.CommittedTotal)
	}
	if got.Remaining != 68628 {
		t.Fatalf("remaining = %s, want £686.28", got.Remaining)
	}
}

func TestBuildJointAccountMatchesSpreadsheetTotals(t *testing.T) {
	joint := BuildJointAccount(nil, []JointAccountItem{
		{Label: "Mortgage", Amount: 119927},
		{Label: "Council Tax", Amount: 17900},
	}, []JointContribution{
		{Member: "Neil", Amount: 100000},
		{Member: "Katie", Amount: 85000},
	}, 2)

	if joint.Total != 137827 {
		t.Fatalf("joint total = %s, want £1378.27", joint.Total)
	}
	if joint.PerPerson != 68914 {
		t.Fatalf("per person = %s, want £689.14", joint.PerPerson)
	}
	if joint.Leftover != 47173 {
		t.Fatalf("leftover = %s, want £471.73", joint.Leftover)
	}
}

func TestBudgetIgnoresMemberNamesAndAllocationLabels(t *testing.T) {
	for _, name := range []string{"Neil", "Katie", "Alex", ""} {
		got := BuildMemberBudget(domain.Member{ID: 1, Name: name}, 10000, []BudgetItem{{Kind: "bill", Amount: 2500}, {Kind: "saving", Amount: 1500}}, []AllocationRule{{Label: "70% - Living Expenses", Amount: 999999}, {Label: "Spendable Income", Amount: 1}})
		if got.SpendableIncome != 6000 || got.Remaining != 6000 {
			t.Fatalf("name/label affected spendable income: %+v", got)
		}
	}
}
func TestJointBudgetUsesHouseholdSize(t *testing.T) {
	for _, n := range []int{0, 1, 3, 5} {
		got := BuildJointAccount(nil, []JointAccountItem{{Amount: 1500}}, nil, n)
		want := domain.Money(0)
		if n > 0 {
			want = 1500 / domain.Money(n)
		}
		if got.PerPerson != want || got.Items[0].PerPerson != want {
			t.Fatalf("%d members: got %d want %d", n, got.PerPerson, want)
		}
	}
}
