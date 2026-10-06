// Package financeconfig loads household financial definitions and resolves their
// configuration references into canonical domain values. It performs no I/O
// beyond explicitly requested file reads and never allocates persisted IDs.
package financeconfig

import (
	"time"

	"github.com/neildavies/swaledale/internal/domain"
)

// Key is a configuration alias, not a database ID or a display name. Aliases are
// unique within their entity kind and scoped to the household being configured.
type Key string

const Version = 1

type Config struct {
	Version    int                  `json:"version"`
	Members    []MemberDefinition   `json:"members"`
	Providers  []ProviderDefinition `json:"providers"`
	Connectors []Key                `json:"connectors"`
	Accounts   []AccountDefinition  `json:"accounts"`
	Bindings   []ProviderBinding    `json:"bindings"`
}

type MemberDefinition struct {
	Key  Key    `json:"key"`
	Name string `json:"name"`
}

type ProviderDefinition struct {
	Key  Key    `json:"key"`
	Name string `json:"name"`
}

// AccountDefinition contains only attributes needed before identity resolution.
// The field types and validation belong to domain, not a second financial model.
type AccountDefinition struct {
	Key      Key                  `json:"key"`
	Name     string               `json:"name"`
	Type     domain.AccountType   `json:"type"`
	Role     domain.FinancialRole `json:"role"`
	Access   domain.AccessClass   `json:"access"`
	Currency domain.Currency      `json:"currency"`
	Owners   []Ownership          `json:"owners"`
}

type Ownership struct {
	Member           Key `json:"member"`
	ShareBasisPoints int `json:"shareBasisPoints"`
}

// ProviderBinding is deliberately incomplete configuration, NOT a validated
// domain.ProviderConnection. Its Key resolves private external identity later.
// Intervals are [ValidFrom, ValidTo); nil end is open. Version 1 allows at most
// one binding at any instant for an account; unbound periods are permitted.
type ProviderBinding struct {
	Key       Key        `json:"key"`
	Account   Key        `json:"account"`
	Provider  Key        `json:"provider"`
	Connector Key        `json:"connector"`
	Product   string     `json:"product,omitempty"`
	ValidFrom time.Time  `json:"validFrom"`
	ValidTo   *time.Time `json:"validTo,omitempty"`
}

func (a AccountDefinition) definition() domain.Account {
	return domain.Account{Name: a.Name, Type: a.Type, Role: a.Role, Access: a.Access, Currency: a.Currency}
}
