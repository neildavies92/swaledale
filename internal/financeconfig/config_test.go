package financeconfig

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/neildavies/swaledale/internal/domain"
)

func fixture() Config {
	return Config{Version: Version,
		Members:    []MemberDefinition{{"member_a", "Alex"}, {"member_b", "Sam"}, {"member_c", "Taylor"}},
		Providers:  []ProviderDefinition{{"bank_a", "Example Bank"}, {"bank_b", "Other Bank"}},
		Connectors: []Key{"manual", "file_import"},
		Accounts:   []AccountDefinition{{Key: "joint_savings", Name: "Joint savings", Type: domain.AccountSavings, Role: domain.RoleJointSavings, Access: domain.AccessAccessible, Currency: domain.GBP, Owners: []Ownership{{"member_a", 7000}, {"member_b", 3000}}}},
		Bindings:   []ProviderBinding{{Key: "savings_initial", Account: "joint_savings", Provider: "bank_a", Connector: "manual", ValidFrom: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}},
	}
}
func identities() Identities {
	return Identities{HouseholdID: 41, Members: map[Key]int64{"member_a": 503, "member_b": 107, "member_c": 902}, Providers: map[Key]domain.ProviderID{"bank_a": 42, "bank_b": 31}, Accounts: map[Key]domain.AccountID{"joint_savings": 805}, Bindings: map[Key]domain.ProviderConnectionID{"savings_initial": 99}}
}
func requireValid(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func encoded(t *testing.T, c Config) string {
	t.Helper()
	data, err := json.Marshal(c)
	requireValid(t, err)
	return string(data)
}

func TestValidAccounts(t *testing.T) {
	cases := []struct {
		name   string
		kind   domain.AccountType
		access domain.AccessClass
		owners []Ownership
	}{
		{"personal current", domain.AccountCurrent, domain.AccessAccessible, []Ownership{{"member_a", 10000}}},
		{"joint equal", domain.AccountSavings, domain.AccessAccessible, []Ownership{{"member_a", 5000}, {"member_b", 5000}}},
		{"joint unequal", domain.AccountSavings, domain.AccessAccessible, []Ownership{{"member_a", 7000}, {"member_b", 3000}}},
		{"three owners", domain.AccountSavings, domain.AccessAccessible, []Ownership{{"member_a", 6000}, {"member_b", 2500}, {"member_c", 1500}}},
		{"credit card", domain.AccountCreditCard, domain.AccessLiability, []Ownership{{"member_b", 10000}}},
		{"ISA", domain.AccountStocksAndSharesISA, domain.AccessAccessible, []Ownership{{"member_a", 10000}}},
		{"SIPP", domain.AccountSIPP, domain.AccessPensionRestricted, []Ownership{{"member_a", 10000}}},
		{"workplace pension", domain.AccountWorkplacePension, domain.AccessPensionRestricted, []Ownership{{"member_c", 10000}}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := fixture()
			c.Accounts[0].Type = tt.kind
			c.Accounts[0].Access = tt.access
			c.Accounts[0].Owners = tt.owners
			loaded, err := Load(strings.NewReader(encoded(t, c)))
			requireValid(t, err)
			resolved, err := loaded.Resolve(identities(), map[Key]string{"savings_initial": "synthetic-account"})
			requireValid(t, err)
			a := resolved.Accounts["joint_savings"]
			if a.ID != 805 || a.HouseholdID != 41 || a.Type != tt.kind {
				t.Fatalf("incorrect resolved account: %+v", a)
			}
			var members []domain.Member
			for _, m := range resolved.Members {
				members = append(members, m)
			}
			requireValid(t, domain.ValidateOwnership(a, resolved.Ownerships["joint_savings"], members))
		})
	}
}
func TestInvalidConfiguration(t *testing.T) {
	cases := []struct {
		name string
		edit func(*Config)
		want string
	}{
		{"version", func(c *Config) { c.Version = 2 }, "version"},
		{"empty members", func(c *Config) { c.Members = nil }, "nonempty"},
		{"empty accounts", func(c *Config) { c.Accounts = nil }, "nonempty"},
		{"missing member", func(c *Config) { c.Accounts[0].Owners[0].Member = "missing" }, "unknown reference"},
		{"missing provider", func(c *Config) { c.Bindings[0].Provider = "missing" }, "unknown reference"},
		{"missing account", func(c *Config) { c.Bindings[0].Account = "missing" }, "unknown reference"},
		{"unknown connector", func(c *Config) { c.Bindings[0].Connector = "not_declared" }, "unknown reference"},
		{"duplicate member", func(c *Config) { c.Members = append(c.Members, c.Members[0]) }, "duplicate key"},
		{"duplicate provider", func(c *Config) { c.Providers = append(c.Providers, c.Providers[0]) }, "duplicate key"},
		{"duplicate account", func(c *Config) { c.Accounts = append(c.Accounts, c.Accounts[0]) }, "duplicate key"},
		{"duplicate connector", func(c *Config) { c.Connectors = append(c.Connectors, c.Connectors[0]) }, "duplicate key"},
		{"duplicate binding", func(c *Config) { c.Bindings = append(c.Bindings, c.Bindings[0]) }, "duplicate key"},
		{"duplicate owner", func(c *Config) { c.Accounts[0].Owners = []Ownership{{"member_a", 5000}, {"member_a", 5000}} }, "duplicate owner"},
		{"bad member key", func(c *Config) { c.Members[0].Key = "Alex Smith" }, "key must match"},
		{"empty account key", func(c *Config) { c.Accounts[0].Key = "" }, "key must match"},
		{"long key", func(c *Config) { c.Providers[0].Key = Key(strings.Repeat("a", 65)) }, "key must match"},
		{"blank member name", func(c *Config) { c.Members[0].Name = " " }, "name: required"},
		{"blank provider name", func(c *Config) { c.Providers[0].Name = " " }, "name: required"},
		{"no owners", func(c *Config) { c.Accounts[0].Owners = nil }, "total 100%"},
		{"low total", func(c *Config) { c.Accounts[0].Owners[0].ShareBasisPoints = 6000 }, "total 100%"},
		{"high total", func(c *Config) { c.Accounts[0].Owners[0].ShareBasisPoints = 8000 }, "exceeds 100%"},
		{"zero share", func(c *Config) { c.Accounts[0].Owners[0].ShareBasisPoints = 0 }, "1..10000"},
		{"negative share", func(c *Config) { c.Accounts[0].Owners[0].ShareBasisPoints = -1 }, "1..10000"},
		{"oversized share", func(c *Config) { c.Accounts[0].Owners[0].ShareBasisPoints = 10001 }, "1..10000"},
		{"type", func(c *Config) { c.Accounts[0].Type = "provider_product" }, "unsupported account type"},
		{"type access", func(c *Config) { c.Accounts[0].Access = domain.AccessLiability }, "requires access"},
		{"currency", func(c *Config) { c.Accounts[0].Currency = "gbp" }, "invalid currency"},
		{"currency length", func(c *Config) { c.Accounts[0].Currency = "" }, "three uppercase"},
		{"blank role", func(c *Config) { c.Accounts[0].Role = " " }, "name and role"},
		{"blank name", func(c *Config) { c.Accounts[0].Name = " " }, "name and role"},
		{"no start", func(c *Config) { c.Bindings[0].ValidFrom = time.Time{} }, "validFrom: required"},
		{"empty interval", func(c *Config) { v := c.Bindings[0].ValidFrom; c.Bindings[0].ValidTo = &v }, "must follow"},
		{"reversed interval", func(c *Config) { v := c.Bindings[0].ValidFrom.Add(-time.Hour); c.Bindings[0].ValidTo = &v }, "must follow"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := fixture()
			tt.edit(&c)
			err := c.Validate()
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("got %v; want %q", err, tt.want)
			}
		})
	}
}
func TestBindingPeriods(t *testing.T) {
	for _, tt := range []struct {
		name      string
		endOffset int
		open      bool
		valid     bool
	}{
		{"adjacent", 0, false, true}, {"gap", -1, false, true}, {"overlap", 1, false, false}, {"open overlap", 0, true, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			c := fixture()
			change := c.Bindings[0].ValidFrom.AddDate(0, 1, 0)
			end := change.Add(time.Duration(tt.endOffset) * time.Hour)
			if !tt.open {
				c.Bindings[0].ValidTo = &end
			}
			next := c.Bindings[0]
			next.Key = "savings_next"
			next.Provider = "bank_b"
			next.ValidFrom = change
			next.ValidTo = nil
			// Reverse input order to prove chronological validation does not depend on file order.
			c.Bindings = append([]ProviderBinding{next}, c.Bindings...)
			before := encoded(t, c)
			err := c.Validate()
			if (err == nil) != tt.valid {
				t.Fatalf("validation = %v", err)
			}
			if encoded(t, c) != before {
				t.Fatal("validation mutated configuration")
			}
		})
	}
	c := fixture()
	c.Bindings = nil
	c.Providers = nil
	c.Connectors = nil
	requireValid(t, c.Validate()) // unbound accounts are valid definitions
}
func TestProviderMigrationKeepsCanonicalIdentityAndHistory(t *testing.T) {
	c := fixture()
	ids := identities()
	external := map[Key]string{"savings_initial": "synthetic-old-account"}
	before, err := c.Resolve(ids, external)
	requireValid(t, err)
	old := before.Connections["savings_initial"]
	snapshot := domain.BalanceSnapshot{ID: 1, AccountID: old.AccountID, Balance: domain.Amount{MinorUnits: 0, Currency: domain.GBP}, AsOf: old.ValidFrom, RecordedAt: old.ValidFrom, Source: domain.ObservationSource{Kind: "manual", ConnectionID: &old.ID}}
	requireValid(t, snapshot.Validate(before.Accounts["joint_savings"]))
	change := old.ValidFrom.AddDate(0, 6, 0)
	c.Bindings[0].ValidTo = &change
	c.Bindings = append(c.Bindings, ProviderBinding{Key: "savings_next", Account: "joint_savings", Provider: "bank_b", Connector: "file_import", Product: "Replacement product", ValidFrom: change})
	ids.Bindings["savings_next"] = 404
	external["savings_next"] = "synthetic-new-account"
	after, err := c.Resolve(ids, external)
	requireValid(t, err)
	if !reflect.DeepEqual(before.Accounts, after.Accounts) || !reflect.DeepEqual(before.Ownerships, after.Ownerships) {
		t.Fatal("provider change altered canonical account or ownership")
	}
	next := after.Connections["savings_next"]
	if next.AccountID != old.AccountID || next.ProviderID == old.ProviderID || next.ID == old.ID || next.ExternalAccountID == old.ExternalAccountID {
		t.Fatal("migration did not replace the mapping independently")
	}
	requireValid(t, snapshot.Validate(after.Accounts["joint_savings"]))
	if *snapshot.Source.ConnectionID != after.Connections["savings_initial"].ID {
		t.Fatal("historical mapping identity lost")
	}
	// Domain results own their end timestamp rather than aliasing mutable config memory.
	expectedEnd := change
	*c.Bindings[0].ValidTo = change.Add(time.Hour)
	if !after.Connections["savings_initial"].ValidTo.Equal(expectedEnd) {
		t.Fatal("resolved mapping timestamp aliases configuration")
	}
}
func TestResolutionIndependentOfOrderNamesAndProviderCode(t *testing.T) {
	c := fixture()
	ids := identities()
	external := map[Key]string{"savings_initial": "synthetic-account"}
	before, err := c.Resolve(ids, external)
	requireValid(t, err)
	c.Members[0], c.Members[2] = c.Members[2], c.Members[0]
	c.Providers[0], c.Providers[1] = c.Providers[1], c.Providers[0]
	after, err := c.Resolve(ids, external)
	requireValid(t, err)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("array order affected resolution")
	}
	c.Members[0].Name = "A different name"
	c.Providers[1].Name = "Entirely new provider brand"
	c.Connectors = append(c.Connectors, "future_connector")
	c.Bindings[0].Connector = "future_connector"
	c.Accounts[0].Role = "future_role"
	after, err = c.Resolve(ids, external)
	requireValid(t, err)
	if after.Connections["savings_initial"].Connector != "future_connector" || after.Accounts["joint_savings"].ID != 805 || !reflect.DeepEqual(before.Ownerships, after.Ownerships) {
		t.Fatal("display names or extensible keys changed ownership/identity")
	}
}
func TestResolutionRejectsIncompleteOrAmbiguousIDs(t *testing.T) {
	cases := []struct {
		name string
		edit func(*Config, *Identities, map[Key]string)
	}{
		{"household", func(_ *Config, i *Identities, _ map[Key]string) { i.HouseholdID = 0 }},
		{"member missing", func(_ *Config, i *Identities, _ map[Key]string) { delete(i.Members, "member_a") }},
		{"member duplicate", func(_ *Config, i *Identities, _ map[Key]string) { i.Members["member_b"] = i.Members["member_a"] }},
		{"provider missing", func(_ *Config, i *Identities, _ map[Key]string) { delete(i.Providers, "bank_a") }},
		{"provider duplicate", func(_ *Config, i *Identities, _ map[Key]string) { i.Providers["bank_b"] = i.Providers["bank_a"] }},
		{"account missing", func(_ *Config, i *Identities, _ map[Key]string) { delete(i.Accounts, "joint_savings") }},
		{"binding missing", func(_ *Config, i *Identities, _ map[Key]string) { delete(i.Bindings, "savings_initial") }},
		{"negative ID", func(_ *Config, i *Identities, _ map[Key]string) { i.Accounts["joint_savings"] = -3 }},
		{"extra alias", func(_ *Config, i *Identities, _ map[Key]string) { i.Accounts["stale"] = 19 }},
		{"private identity missing", func(_ *Config, _ *Identities, e map[Key]string) { delete(e, "savings_initial") }},
		{"private identity blank", func(_ *Config, _ *Identities, e map[Key]string) { e["savings_initial"] = " " }},
		{"private identity extra", func(_ *Config, _ *Identities, e map[Key]string) {
			e["stale"] = "private-value-must-not-appear-in-errors"
		}},
		{"invalid config", func(c *Config, _ *Identities, _ map[Key]string) { c.Accounts[0].Owners = nil }},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			c := fixture()
			ids := identities()
			external := map[Key]string{"savings_initial": "private-value-must-not-appear-in-errors"}
			tt.edit(&c, &ids, external)
			result, err := c.Resolve(ids, external)
			if err == nil {
				t.Fatal("expected error")
			}
			if !reflect.DeepEqual(result, Resolved{}) {
				t.Fatal("returned partial financial records on failure")
			}
			if strings.Contains(err.Error(), "private-value-must-not-appear-in-errors") {
				t.Fatal("error exposed private identity")
			}
		})
	}
}

func TestMultipleAccountsResolveByKey(t *testing.T) {
	c := fixture()
	ids := identities()
	external := map[Key]string{"savings_initial": "synthetic-a"}
	second := c.Accounts[0]
	second.Key = "second_account"
	second.Name = "Second account"
	c.Accounts = append(c.Accounts, second)
	ids.Accounts[second.Key] = 302
	binding := c.Bindings[0]
	binding.Key = "second_binding"
	binding.Account = second.Key
	c.Bindings = append(c.Bindings, binding)
	ids.Bindings[binding.Key] = 201
	external[binding.Key] = "synthetic-b"
	before, err := c.Resolve(ids, external)
	requireValid(t, err)
	c.Accounts[0], c.Accounts[1] = c.Accounts[1], c.Accounts[0]
	c.Bindings[0], c.Bindings[1] = c.Bindings[1], c.Bindings[0]
	after, err := c.Resolve(ids, external)
	requireValid(t, err)
	if !reflect.DeepEqual(before, after) {
		t.Fatal("account/binding array order affected identity")
	}
	ids.Accounts[second.Key] = 805
	if _, err := c.Resolve(ids, external); err == nil {
		t.Fatal("duplicate resolved account ID accepted")
	}
	ids.Accounts[second.Key] = 302
	ids.Bindings[binding.Key] = 99
	if _, err := c.Resolve(ids, external); err == nil {
		t.Fatal("duplicate resolved binding ID accepted")
	}
}
