package domain

import (
	"testing"
	"time"
)

func testMembers() []Member {
	return []Member{{ID: 1, HouseholdID: 7, Name: "Alex"}, {ID: 2, HouseholdID: 7, Name: "Sam"}, {ID: 3, HouseholdID: 7, Name: "Taylor"}}
}
func testAccount(kind AccountType, access AccessClass) Account {
	return Account{ID: 10, HouseholdID: 7, Name: "Household account", Type: kind, Role: RoleSalary, Access: access, Currency: GBP}
}
func mustValid(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func mustInvalid(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestAccountScenarios(t *testing.T) {
	cases := []struct {
		name   string
		kind   AccountType
		access AccessClass
		role   FinancialRole
		shares []int
	}{
		{"personal current", AccountCurrent, AccessAccessible, RoleSalary, []int{10000}},
		{"joint savings equal", AccountSavings, AccessAccessible, RoleJointSavings, []int{5000, 5000}},
		{"joint savings unequal", AccountSavings, AccessAccessible, RoleJointSavings, []int{7000, 3000}},
		{"three owners", AccountSavings, AccessAccessible, RoleJointSavings, []int{6000, 2500, 1500}},
		{"credit card liability", AccountCreditCard, AccessLiability, RoleSpending, []int{10000}},
		{"stocks and shares ISA", AccountStocksAndSharesISA, AccessAccessible, RoleISABridge, []int{10000}},
		{"SIPP", AccountSIPP, AccessPensionRestricted, RolePension, []int{10000}},
		{"workplace pension", AccountWorkplacePension, AccessPensionRestricted, RolePension, []int{10000}},
		{"investment", AccountInvestment, AccessAccessible, RoleLongTermInvestment, []int{10000}},
		{"new role", AccountSavings, AccessAccessible, FinancialRole("education_savings"), []int{10000}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			a := testAccount(tt.kind, tt.access)
			a.Role = tt.role
			owners := make([]AccountOwnership, len(tt.shares))
			for i, share := range tt.shares {
				owners[i] = AccountOwnership{a.ID, int64(i + 1), share}
			}
			mustValid(t, ValidateOwnership(a, owners, testMembers()))
			members := testMembers()
			members[0].Name = "Neil"
			members[1].Name = "70% - Living Expenses"
			mustValid(t, ValidateOwnership(a, owners, members))
		})
	}
}

func TestOwnershipRejectsInvalidSets(t *testing.T) {
	a := testAccount(AccountSavings, AccessAccessible)
	cases := map[string][]AccountOwnership{
		"empty":           nil,
		"zero share":      {{10, 1, 0}, {10, 2, 10000}},
		"negative share":  {{10, 1, -1}, {10, 2, 10001}},
		"over limit":      {{10, 1, 10001}},
		"underallocated":  {{10, 1, 9999}},
		"overallocated":   {{10, 1, 6000}, {10, 2, 6000}},
		"duplicate owner": {{10, 1, 5000}, {10, 1, 5000}},
		"unknown member":  {{10, 4, 10000}},
		"wrong account":   {{11, 1, 10000}},
		"zero member":     {{10, 0, 10000}},
	}
	for name, owners := range cases {
		t.Run(name, func(t *testing.T) { mustInvalid(t, ValidateOwnership(a, owners, testMembers())) })
	}
	members := testMembers()
	members[0].HouseholdID = 8
	mustInvalid(t, ValidateOwnership(a, []AccountOwnership{{10, 1, 10000}}, members))
	mustInvalid(t, ValidateOwnership(a, []AccountOwnership{{10, 1, 10000}}, append(testMembers(), testMembers()[0])))
	members = testMembers()
	members[0].ID = 0
	mustInvalid(t, ValidateOwnership(a, []AccountOwnership{{10, 1, 10000}}, members))
}

func TestAccountValidation(t *testing.T) {
	for name, mutate := range map[string]func(*Account){
		"zero ID": func(a *Account) { a.ID = 0 }, "zero household": func(a *Account) { a.HouseholdID = 0 },
		"empty name": func(a *Account) { a.Name = " " }, "empty role": func(a *Account) { a.Role = "" },
		"unknown type": func(a *Account) { a.Type = "provider_product" }, "missing access": func(a *Account) { a.Access = "" },
		"cash pension mismatch":       func(a *Account) { a.Access = AccessPensionRestricted },
		"card accessible mismatch":    func(a *Account) { a.Type = AccountCreditCard },
		"pension accessible mismatch": func(a *Account) { a.Type = AccountSIPP },
		"invalid currency":            func(a *Account) { a.Currency = "gbp" },
	} {
		t.Run(name, func(t *testing.T) {
			a := testAccount(AccountCurrent, AccessAccessible)
			mutate(&a)
			mustInvalid(t, a.Validate())
		})
	}
}

func TestProviderChangePreservesAccountAndHistory(t *testing.T) {
	a := testAccount(AccountStocksAndSharesISA, AccessAccessible)
	original := a
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	change := start.AddDate(0, 6, 0)
	providers := []Provider{{ID: 1, Name: "Vanguard"}, {ID: 2, Name: "Another platform"}}
	old := ProviderConnection{ID: 100, AccountID: a.ID, ProviderID: 1, Connector: "csv", ExternalAccountID: "external-123", Product: "ISA", ValidFrom: start, ValidTo: &change}
	next := ProviderConnection{ID: 101, AccountID: a.ID, ProviderID: 2, Connector: "api", ExternalAccountID: "external-123", Product: "New ISA product", ValidFrom: change}
	mustValid(t, old.Validate(a, providers[0]))
	mustValid(t, next.Validate(a, providers[1]))
	history := BalanceSnapshot{ID: 1, AccountID: a.ID, Balance: Amount{50000, GBP}, AsOf: start, RecordedAt: start, Source: ObservationSource{Kind: "csv", ConnectionID: &old.ID, ExternalID: "row-1"}}
	mustValid(t, history.Validate(a))
	if a != original || *history.Source.ConnectionID != old.ID {
		t.Fatal("mapping replacement changed identity/history")
	}
	for name, mutate := range map[string]func(*ProviderConnection){
		"identity": func(p *ProviderConnection) { p.ID = 0 }, "account": func(p *ProviderConnection) { p.AccountID = 999 },
		"provider": func(p *ProviderConnection) { p.ProviderID = 2 }, "external": func(p *ProviderConnection) { p.ExternalAccountID = " " },
		"connector": func(p *ProviderConnection) { p.Connector = "" }, "start": func(p *ProviderConnection) { p.ValidFrom = time.Time{} },
		"end": func(p *ProviderConnection) { p.ValidTo = &start },
	} {
		t.Run(name, func(t *testing.T) { p := old; mutate(&p); mustInvalid(t, p.Validate(a, providers[0])) })
	}
}

func TestDefinitionValidationDoesNotReplaceIdentityValidation(t *testing.T) {
	a := testAccount(AccountSavings, AccessAccessible)
	a.ID = 0
	a.HouseholdID = 0
	mustValid(t, a.ValidateDefinition())
	mustInvalid(t, a.Validate())
	a.Access = AccessLiability
	mustInvalid(t, a.ValidateDefinition())
}
