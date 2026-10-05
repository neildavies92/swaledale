package domain

import (
	"fmt"
	"strings"
	"time"
)

type AccountID int64
type ProviderID int64
type ProviderConnectionID int64
type InstrumentID int64

type AccountType string

const (
	AccountCurrent            AccountType = "current"
	AccountSavings            AccountType = "savings"
	AccountCreditCard         AccountType = "credit_card"
	AccountInvestment         AccountType = "investment"
	AccountStocksAndSharesISA AccountType = "stocks_and_shares_isa"
	AccountSIPP               AccountType = "sipp"
	AccountWorkplacePension   AccountType = "workplace_pension"
)

func (t AccountType) InvestmentCapable() bool {
	return t == AccountInvestment || t == AccountStocksAndSharesISA || t == AccountSIPP || t == AccountWorkplacePension
}

type AccessClass string

const (
	AccessAccessible        AccessClass = "accessible"
	AccessPensionRestricted AccessClass = "pension_restricted"
	AccessLiability         AccessClass = "liability"
)

// FinancialRole is an extensible configuration key, not a provider name. Common
// roles are conveniences; new nonblank keys do not require changing an enum.
type FinancialRole string

const (
	RoleSalary             FinancialRole = "salary"
	RoleSpending           FinancialRole = "spending"
	RoleEmergencyFund      FinancialRole = "emergency_fund"
	RoleJointBills         FinancialRole = "joint_bills"
	RoleJointSavings       FinancialRole = "joint_savings"
	RoleLongTermInvestment FinancialRole = "long_term_investment"
	RoleISABridge          FinancialRole = "isa_bridge"
	RolePension            FinancialRole = "pension"
)

// Account is stable identity; neither provider nor current balance nor instrument
// is embedded in it. Household/member IDs keep the existing store's int64 contract.
type Account struct {
	ID          AccountID     `json:"id"`
	HouseholdID int64         `json:"householdId"`
	Name        string        `json:"name"`
	Type        AccountType   `json:"type"`
	Role        FinancialRole `json:"role"`
	Access      AccessClass   `json:"access"`
	Currency    Currency      `json:"currency"`
}

func (a Account) Validate() error {
	if a.ID <= 0 || a.HouseholdID <= 0 || strings.TrimSpace(a.Name) == "" || strings.TrimSpace(string(a.Role)) == "" {
		return fmt.Errorf("account requires positive identity, household, name and role")
	}
	if err := a.Currency.Validate(); err != nil {
		return err
	}
	expected := AccessAccessible
	switch a.Type {
	case AccountCurrent, AccountSavings, AccountInvestment, AccountStocksAndSharesISA:
	case AccountCreditCard:
		expected = AccessLiability
	case AccountSIPP, AccountWorkplacePension:
		expected = AccessPensionRestricted
	default:
		return fmt.Errorf("unsupported account type %q", a.Type)
	}
	if a.Access != expected {
		return fmt.Errorf("account type %s requires access %s", a.Type, expected)
	}
	return nil
}

const FullOwnership = 10000

// AccountOwnership is one member's share in basis points (1 = 0.01%). A complete
// set totals 10,000. Version sets in persistence before supporting ownership edits.
type AccountOwnership struct {
	AccountID        AccountID `json:"accountId"`
	MemberID         int64     `json:"memberId"`
	ShareBasisPoints int       `json:"shareBasisPoints"`
}

// ValidateOwnership validates a complete set against known household members.
// The repository must supply authoritative members; this is not authorization.
func ValidateOwnership(account Account, owners []AccountOwnership, members []Member) error {
	if err := account.Validate(); err != nil {
		return err
	}
	known := make(map[int64]int64, len(members))
	for _, m := range members {
		if m.ID <= 0 || m.HouseholdID <= 0 {
			return fmt.Errorf("invalid member identity")
		}
		if _, exists := known[m.ID]; exists {
			return fmt.Errorf("duplicate member %d", m.ID)
		}
		known[m.ID] = m.HouseholdID
	}
	total := 0
	seen := make(map[int64]bool, len(owners))
	for _, o := range owners {
		if o.AccountID != account.ID || o.MemberID <= 0 || known[o.MemberID] != account.HouseholdID {
			return fmt.Errorf("ownership must reference this account and a member of its household")
		}
		if seen[o.MemberID] || o.ShareBasisPoints <= 0 || o.ShareBasisPoints > FullOwnership {
			return fmt.Errorf("ownership requires unique members and shares in 1..10000")
		}
		seen[o.MemberID] = true
		total += o.ShareBasisPoints
		if total > FullOwnership {
			return fmt.Errorf("ownership exceeds 100%%")
		}
	}
	if total != FullOwnership {
		return fmt.Errorf("ownership must total 100%%")
	}
	return nil
}

type Provider struct {
	ID   ProviderID `json:"id"`
	Name string     `json:"name"`
}

// ProviderConnection is an effective-dated mapping of one stable account to an
// external provider account and connector. It contains no credentials. ValidTo is
// exclusive; nil means open-ended. A migration closes one mapping and adds another.
type ProviderConnection struct {
	ID                ProviderConnectionID `json:"id"`
	AccountID         AccountID            `json:"accountId"`
	ProviderID        ProviderID           `json:"providerId"`
	Connector         string               `json:"connector"`
	ExternalAccountID string               `json:"externalAccountId"`
	Product           string               `json:"product,omitempty"`
	ValidFrom         time.Time            `json:"validFrom"`
	ValidTo           *time.Time           `json:"validTo,omitempty"`
}

func (p ProviderConnection) Validate(account Account, provider Provider) error {
	if err := account.Validate(); err != nil {
		return err
	}
	if p.ID <= 0 || p.AccountID != account.ID || provider.ID <= 0 || p.ProviderID != provider.ID || strings.TrimSpace(provider.Name) == "" {
		return fmt.Errorf("provider mapping requires valid internal account and provider identities")
	}
	if strings.TrimSpace(p.Connector) == "" || strings.TrimSpace(p.ExternalAccountID) == "" || p.ValidFrom.IsZero() {
		return fmt.Errorf("provider mapping requires connector, external account and start time")
	}
	if p.ValidTo != nil && !p.ValidTo.After(p.ValidFrom) {
		return fmt.Errorf("provider mapping end must follow start")
	}
	return nil
}

// ObservationSource identifies manual/import provenance or a retained mapping.
// ExternalID is optional and meaningful only within its source namespace.
type ObservationSource struct {
	Kind         string                `json:"kind"`
	ConnectionID *ProviderConnectionID `json:"connectionId,omitempty"`
	ExternalID   string                `json:"externalId,omitempty"`
}

func (s ObservationSource) Validate() error {
	if strings.TrimSpace(s.Kind) == "" {
		return fmt.Errorf("observation source is required")
	}
	if s.ConnectionID != nil && *s.ConnectionID <= 0 {
		return fmt.Errorf("invalid source connection")
	}
	return nil
}
