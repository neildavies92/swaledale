package financeconfig

import (
	"fmt"
	"sort"
	"strings"

	"github.com/neildavies/swaledale/internal/domain"
)

// Identities must be supplied by a trusted persistence/bootstrap boundary. Map
// entries are persisted aliases, never values generated from configuration order.
// The caller is responsible for verifying their existence and household scope.
type Identities struct {
	HouseholdID int64
	Members     map[Key]int64
	Providers   map[Key]domain.ProviderID
	Accounts    map[Key]domain.AccountID
	Bindings    map[Key]domain.ProviderConnectionID
}

// Resolved contains canonical values keyed by their configuration aliases.
// Connections contain private external IDs: do not log or commit this result.
type Resolved struct {
	Members     map[Key]domain.Member
	Providers   map[Key]domain.Provider
	Accounts    map[Key]domain.Account
	Ownerships  map[Key][]domain.AccountOwnership
	Connections map[Key]domain.ProviderConnection
}

// Resolve is a pure conversion, not persistence or a migration operation.
// externalAccounts is supplied privately at runtime, keyed by binding key. No
// credentials are needed. Missing IDs fail closed; no placeholders are invented.
func (c Config) Resolve(ids Identities, externalAccounts map[Key]string) (Resolved, error) {
	if err := c.Validate(); err != nil {
		return Resolved{}, err
	}
	if ids.HouseholdID <= 0 {
		return Resolved{}, fmt.Errorf("resolved household ID must be positive")
	}
	memberKeys := make([]Key, 0, len(c.Members))
	for _, m := range c.Members {
		memberKeys = append(memberKeys, m.Key)
	}
	providerKeys := make([]Key, 0, len(c.Providers))
	for _, p := range c.Providers {
		providerKeys = append(providerKeys, p.Key)
	}
	accountKeys := make([]Key, 0, len(c.Accounts))
	for _, a := range c.Accounts {
		accountKeys = append(accountKeys, a.Key)
	}
	bindingKeys := make([]Key, 0, len(c.Bindings))
	for _, b := range c.Bindings {
		bindingKeys = append(bindingKeys, b.Key)
	}
	if err := validateIDs(memberKeys, ids.Members, "members"); err != nil {
		return Resolved{}, err
	}
	if err := validateIDs(providerKeys, ids.Providers, "providers"); err != nil {
		return Resolved{}, err
	}
	if err := validateIDs(accountKeys, ids.Accounts, "accounts"); err != nil {
		return Resolved{}, err
	}
	if err := validateIDs(bindingKeys, ids.Bindings, "bindings"); err != nil {
		return Resolved{}, err
	}
	for _, key := range bindingKeys {
		if strings.TrimSpace(externalAccounts[key]) == "" {
			return Resolved{}, fmt.Errorf("binding %q: private external identity is required", key)
		}
	}
	if err := exactKeys(bindingKeys, externalAccounts, "private external identities"); err != nil {
		return Resolved{}, err
	}
	result := Resolved{
		Members: map[Key]domain.Member{}, Providers: map[Key]domain.Provider{}, Accounts: map[Key]domain.Account{},
		Ownerships: map[Key][]domain.AccountOwnership{}, Connections: map[Key]domain.ProviderConnection{},
	}
	members := make([]domain.Member, 0, len(c.Members))
	for _, m := range c.Members {
		resolved := domain.Member{ID: ids.Members[m.Key], HouseholdID: ids.HouseholdID, Name: m.Name}
		result.Members[m.Key] = resolved
		members = append(members, resolved)
	}
	for _, p := range c.Providers {
		result.Providers[p.Key] = domain.Provider{ID: ids.Providers[p.Key], Name: p.Name}
	}
	for _, a := range c.Accounts {
		resolved := a.definition()
		resolved.ID = ids.Accounts[a.Key]
		resolved.HouseholdID = ids.HouseholdID
		owners := make([]domain.AccountOwnership, 0, len(a.Owners))
		for _, o := range a.Owners {
			owners = append(owners, domain.AccountOwnership{AccountID: resolved.ID, MemberID: ids.Members[o.Member], ShareBasisPoints: o.ShareBasisPoints})
		}
		if err := domain.ValidateOwnership(resolved, owners, members); err != nil {
			return Resolved{}, fmt.Errorf("account %q: %w", a.Key, err)
		}
		result.Accounts[a.Key] = resolved
		result.Ownerships[a.Key] = owners
	}
	for _, b := range c.Bindings {
		connection := domain.ProviderConnection{
			ID: ids.Bindings[b.Key], AccountID: ids.Accounts[b.Account], ProviderID: ids.Providers[b.Provider],
			Connector: string(b.Connector), ExternalAccountID: externalAccounts[b.Key], Product: b.Product, ValidFrom: b.ValidFrom,
		}
		if b.ValidTo != nil {
			end := *b.ValidTo
			connection.ValidTo = &end
		}
		if err := connection.Validate(result.Accounts[b.Account], result.Providers[b.Provider]); err != nil {
			return Resolved{}, fmt.Errorf("binding %q: %w", b.Key, err)
		}
		result.Connections[b.Key] = connection
	}
	return result, nil
}

func validateIDs[T ~int64](keys []Key, ids map[Key]T, kind string) error {
	seen := map[T]bool{}
	for _, key := range keys {
		id := ids[key]
		if id <= 0 {
			return fmt.Errorf("%s %q: positive resolved ID required", kind, key)
		}
		if seen[id] {
			return fmt.Errorf("%s %q: resolved ID is already assigned to another key", kind, key)
		}
		seen[id] = true
	}
	return exactKeys(keys, ids, kind)
}
func exactKeys[T any](keys []Key, values map[Key]T, kind string) error {
	expected := map[Key]bool{}
	for _, key := range keys {
		expected[key] = true
	}
	extras := []string{}
	for key := range values {
		if !expected[key] {
			extras = append(extras, string(key))
		}
	}
	sort.Strings(extras)
	if len(extras) > 0 {
		return fmt.Errorf("%s: unexpected key %q", kind, extras[0])
	}
	return nil
}
