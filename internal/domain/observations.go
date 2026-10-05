package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// BalanceSnapshot is a signed net balance: assets positive, debt negative.
// A card overpayment may be positive. AccessLiability identifies the account;
// it must not cause a second negation of the amount in later net-worth calculations.
type BalanceSnapshot struct {
	ID         int64             `json:"id"`
	AccountID  AccountID         `json:"accountId"`
	Balance    Amount            `json:"balance"`
	AsOf       time.Time         `json:"asOf"`
	RecordedAt time.Time         `json:"recordedAt"`
	Source     ObservationSource `json:"source"`
}

func validateObservation(id int64, accountID AccountID, account Account, amount Amount, asOf, recordedAt time.Time, source ObservationSource) error {
	if err := account.Validate(); err != nil {
		return err
	}
	if id <= 0 || accountID != account.ID {
		return fmt.Errorf("observation must belong to account")
	}
	if err := amount.Validate(); err != nil {
		return err
	}
	if amount.Currency != account.Currency {
		return fmt.Errorf("observation currency must match account reporting currency")
	}
	if asOf.IsZero() || recordedAt.IsZero() || recordedAt.Before(asOf) {
		return fmt.Errorf("observation requires as-of and recording times in order")
	}
	return source.Validate()
}
func (b BalanceSnapshot) Validate(account Account) error {
	return validateObservation(b.ID, b.AccountID, account, b.Balance, b.AsOf, b.RecordedAt, b.Source)
}

// Instrument identity is independent of the account/platform. Identifiers are
// optional namespaced data (e.g. isin, exchange:ticker, provider-specific keys).
type Instrument struct {
	ID          InstrumentID      `json:"id"`
	Name        string            `json:"name"`
	Identifiers map[string]string `json:"identifiers,omitempty"`
}

func (i Instrument) Validate() error {
	if i.ID <= 0 || strings.TrimSpace(i.Name) == "" {
		return fmt.Errorf("instrument requires identity and name")
	}
	for namespace, value := range i.Identifiers {
		if strings.TrimSpace(namespace) == "" || strings.TrimSpace(value) == "" {
			return fmt.Errorf("instrument identifiers require namespace and value")
		}
	}
	return nil
}

// Quantity is an exact nonnegative decimal string, never a binary float. No
// arithmetic/precision limit is imposed here; persistence should use NUMERIC.
type Quantity string

var quantityPattern = regexp.MustCompile(`^(0|[1-9][0-9]*)(\.[0-9]+)?$`)

func (q Quantity) Validate() error {
	if !quantityPattern.MatchString(string(q)) {
		return fmt.Errorf("quantity must be a nonnegative decimal string")
	}
	return nil
}

type Holding struct {
	ID           int64             `json:"id"`
	AccountID    AccountID         `json:"accountId"`
	InstrumentID InstrumentID      `json:"instrumentId"`
	Quantity     Quantity          `json:"quantity"`
	Valuation    Amount            `json:"valuation"`
	AsOf         time.Time         `json:"asOf"`
	RecordedAt   time.Time         `json:"recordedAt"`
	Source       ObservationSource `json:"source"`
}

func (h Holding) Validate(account Account, instrument Instrument) error {
	if err := validateObservation(h.ID, h.AccountID, account, h.Valuation, h.AsOf, h.RecordedAt, h.Source); err != nil {
		return err
	}
	if !account.Type.InvestmentCapable() {
		return fmt.Errorf("holdings require an investment-capable account")
	}
	if err := instrument.Validate(); err != nil {
		return err
	}
	if h.InstrumentID != instrument.ID {
		return fmt.Errorf("holding must reference instrument")
	}
	if h.Valuation.MinorUnits < 0 {
		return fmt.Errorf("holding valuation cannot be negative")
	}
	return h.Quantity.Validate()
}

// Transaction records booked entries only. Amount is signed: inflow positive,
// outflow negative (a card purchase is negative, a repayment positive). Pending
// reconciliation, categorisation and transfer matching belong to later tickets.
type Transaction struct {
	ID          int64             `json:"id"`
	AccountID   AccountID         `json:"accountId"`
	Amount      Amount            `json:"amount"`
	BookedAt    time.Time         `json:"bookedAt"`
	RecordedAt  time.Time         `json:"recordedAt"`
	Description string            `json:"description"`
	Source      ObservationSource `json:"source"`
}

func (t Transaction) Validate(account Account) error {
	return validateObservation(t.ID, t.AccountID, account, t.Amount, t.BookedAt, t.RecordedAt, t.Source)
}

// FIRESnapshot is derived output, not a calculation engine. Nil MemberID denotes
// the household. Optional FIRENumber/ProgressBasisPoints distinguish uncomputed
// metrics from zero. All amounts use Currency; liabilities are positive debt totals.
// NetWorth is signed; ProgressBasisPoints may exceed 10000.
type FIRESnapshot struct {
	ID                  int64     `json:"id"`
	HouseholdID         int64     `json:"householdId"`
	MemberID            *int64    `json:"memberId,omitempty"`
	AsOf                time.Time `json:"asOf"`
	CalculatedAt        time.Time `json:"calculatedAt"`
	Currency            Currency  `json:"currency"`
	AccessibleAssets    Amount    `json:"accessibleAssets"`
	PensionAssets       Amount    `json:"pensionAssets"`
	Liabilities         Amount    `json:"liabilities"`
	NetWorth            Amount    `json:"netWorth"`
	FIRENumber          *Amount   `json:"fireNumber,omitempty"`
	ProgressBasisPoints *int64    `json:"progressBasisPoints,omitempty"`
	CalculationVersion  string    `json:"calculationVersion"`
	AssumptionsRef      string    `json:"assumptionsRef"`
}

func (f FIRESnapshot) Validate(household Household, members []Member) error {
	if f.ID <= 0 || household.ID <= 0 || f.HouseholdID != household.ID {
		return fmt.Errorf("FIRE snapshot requires household identity")
	}
	if err := f.Currency.Validate(); err != nil {
		return err
	}
	if f.MemberID != nil {
		found := false
		for _, m := range members {
			if m.ID > 0 && m.ID == *f.MemberID && m.HouseholdID == f.HouseholdID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("FIRE member must belong to household")
		}
	}
	if f.AsOf.IsZero() || f.CalculatedAt.IsZero() || f.CalculatedAt.Before(f.AsOf) || strings.TrimSpace(f.CalculationVersion) == "" || strings.TrimSpace(f.AssumptionsRef) == "" {
		return fmt.Errorf("FIRE snapshot requires times, calculation version and assumptions reference")
	}
	for _, amount := range []Amount{f.AccessibleAssets, f.PensionAssets, f.Liabilities, f.NetWorth} {
		if amount.Currency != f.Currency {
			return fmt.Errorf("FIRE metrics must share reporting currency")
		}
	}
	if f.AccessibleAssets.MinorUnits < 0 || f.PensionAssets.MinorUnits < 0 || f.Liabilities.MinorUnits < 0 {
		return fmt.Errorf("asset and debt totals cannot be negative")
	}
	if f.FIRENumber != nil && (f.FIRENumber.Currency != f.Currency || f.FIRENumber.MinorUnits <= 0) {
		return fmt.Errorf("FIRE target must be positive in reporting currency")
	}
	if f.ProgressBasisPoints != nil && (f.FIRENumber == nil || *f.ProgressBasisPoints < 0) {
		return fmt.Errorf("FIRE progress requires a target and nonnegative basis points")
	}
	return nil
}
