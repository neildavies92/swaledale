package domain

import (
	"encoding/json"
	"testing"
	"time"
)

var observed = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func testBalance() BalanceSnapshot {
	return BalanceSnapshot{ID: 1, AccountID: 10, Balance: Amount{12345, GBP}, AsOf: observed, RecordedAt: observed.Add(time.Hour), Source: ObservationSource{Kind: "manual"}}
}
func TestHistoricalBalances(t *testing.T) {
	a := testAccount(AccountSavings, AccessAccessible)
	first := testBalance()
	second := first
	second.ID = 2
	second.AsOf = observed.AddDate(0, 1, 0)
	second.RecordedAt = second.AsOf.Add(time.Hour)
	second.Balance.MinorUnits = 15000
	for _, b := range []BalanceSnapshot{first, second} {
		mustValid(t, b.Validate(a))
		data, err := json.Marshal(b)
		mustValid(t, err)
		var restored BalanceSnapshot
		mustValid(t, json.Unmarshal(data, &restored))
		if !restored.AsOf.Equal(b.AsOf) || !restored.RecordedAt.Equal(b.RecordedAt) || restored.Balance != b.Balance {
			t.Fatal("snapshot history lost in serialization")
		}
	}
	if first.Balance.MinorUnits != 12345 || !first.AsOf.Equal(observed) {
		t.Fatal("later snapshot changed history")
	}
	card := testAccount(AccountCreditCard, AccessLiability)
	debt := testBalance()
	debt.Balance.MinorUnits = -10000
	mustValid(t, debt.Validate(card))
	debt.Balance.MinorUnits = 100
	mustValid(t, debt.Validate(card)) // overpayment
	cash := testAccount(AccountCurrent, AccessAccessible)
	debt.Balance.MinorUnits = -100
	mustValid(t, debt.Validate(cash)) // overdraft
}
func TestInvalidObservations(t *testing.T) {
	a := testAccount(AccountSavings, AccessAccessible)
	for name, mutate := range map[string]func(*BalanceSnapshot){
		"zero ID": func(b *BalanceSnapshot) { b.ID = 0 }, "wrong account": func(b *BalanceSnapshot) { b.AccountID = 11 },
		"currency": func(b *BalanceSnapshot) { b.Balance.Currency = "EUR" }, "invalid currency": func(b *BalanceSnapshot) { b.Balance.Currency = "" },
		"missing as of": func(b *BalanceSnapshot) { b.AsOf = time.Time{} }, "missing recorded": func(b *BalanceSnapshot) { b.RecordedAt = time.Time{} },
		"reversed times": func(b *BalanceSnapshot) { b.RecordedAt = b.AsOf.Add(-time.Hour) },
		"source":         func(b *BalanceSnapshot) { b.Source.Kind = "" }, "connection": func(b *BalanceSnapshot) { id := ProviderConnectionID(0); b.Source.ConnectionID = &id },
	} {
		t.Run(name, func(t *testing.T) { b := testBalance(); mutate(&b); mustInvalid(t, b.Validate(a)) })
	}
}
func TestHoldingsIndependentOfAccount(t *testing.T) {
	vall := Instrument{ID: 20, Name: "VALL", Identifiers: map[string]string{"ticker": "VALL"}}
	replacement := Instrument{ID: 21, Name: "Another fund"} // no ticker/ISIN required
	for _, kind := range []AccountType{AccountStocksAndSharesISA, AccountSIPP, AccountWorkplacePension} {
		access := AccessAccessible
		if kind != AccountStocksAndSharesISA {
			access = AccessPensionRestricted
		}
		a := testAccount(kind, access)
		original := a
		h := Holding{ID: 1, AccountID: a.ID, InstrumentID: vall.ID, Quantity: "123.456789123", Valuation: Amount{98765, GBP}, AsOf: observed, RecordedAt: observed, Source: ObservationSource{Kind: "import"}}
		mustValid(t, h.Validate(a, vall))
		next := h
		next.ID = 2
		next.InstrumentID = replacement.ID
		next.AsOf = observed.Add(time.Hour)
		next.RecordedAt = next.AsOf
		mustValid(t, next.Validate(a, replacement))
		if a != original || h.InstrumentID != vall.ID {
			t.Fatal("holding change changed account or past holding")
		}
		data, err := json.Marshal(h)
		mustValid(t, err)
		var restored Holding
		mustValid(t, json.Unmarshal(data, &restored))
		if restored.Quantity != h.Quantity {
			t.Fatal("quantity lost precision")
		}
		mustInvalid(t, h.Validate(testAccount(AccountCurrent, AccessAccessible), vall))
		mustInvalid(t, h.Validate(a, replacement))
		h.Valuation.MinorUnits = -1
		mustInvalid(t, h.Validate(a, vall))
	}
	for _, q := range []Quantity{"0", "0.0000000001", "999999999999999999999.123456789", "1"} {
		mustValid(t, q.Validate())
	}
	for _, q := range []Quantity{"", "-1", "NaN", "1e3", "1/3", "1.", ".5", "01", "+1"} {
		mustInvalid(t, q.Validate())
	}
	mustInvalid(t, (Instrument{ID: 1}).Validate())
	mustInvalid(t, (Instrument{Name: "fund"}).Validate())
	mustInvalid(t, (Instrument{ID: 1, Name: "fund", Identifiers: map[string]string{"isin": ""}}).Validate())
}
func TestTransactions(t *testing.T) {
	a := testAccount(AccountCurrent, AccessAccessible)
	tx := Transaction{ID: 1, AccountID: a.ID, Amount: Amount{-2500, GBP}, BookedAt: observed, RecordedAt: observed, Description: "Purchase", Source: ObservationSource{Kind: "import", ExternalID: "provider-entry-1"}}
	mustValid(t, tx.Validate(a))
	tx.Amount.MinorUnits = 100000
	mustValid(t, tx.Validate(a))
	tx.AccountID = 99
	mustInvalid(t, tx.Validate(a))
}
func testFIRE() FIRESnapshot {
	return FIRESnapshot{ID: 1, HouseholdID: 7, AsOf: observed, CalculatedAt: observed.Add(time.Hour), Currency: GBP,
		AccessibleAssets: Amount{100000, GBP}, PensionAssets: Amount{200000, GBP}, Liabilities: Amount{50000, GBP}, NetWorth: Amount{250000, GBP}, CalculationVersion: "v1", AssumptionsRef: "assumptions-1"}
}
func TestFIREScopesAndOptionalMetrics(t *testing.T) {
	household := Household{ID: 7, Name: "Household", Currency: "GBP"}
	f := testFIRE()
	mustValid(t, f.Validate(household, testMembers()))
	id := int64(3)
	f.MemberID = &id
	mustValid(t, f.Validate(household, testMembers()))
	target := Amount{1000000, GBP}
	progress := int64(12000)
	f.FIRENumber = &target
	f.ProgressBasisPoints = &progress
	mustValid(t, f.Validate(household, testMembers()))
	for name, mutate := range map[string]func(*FIRESnapshot){
		"identity": func(f *FIRESnapshot) { f.ID = 0 }, "household": func(f *FIRESnapshot) { f.HouseholdID = 8 },
		"scope": func(f *FIRESnapshot) { v := int64(999); f.MemberID = &v }, "timestamp": func(f *FIRESnapshot) { f.AsOf = time.Time{} },
		"recording": func(f *FIRESnapshot) { f.CalculatedAt = f.AsOf.Add(-time.Hour) }, "version": func(f *FIRESnapshot) { f.CalculationVersion = "" },
		"assumptions": func(f *FIRESnapshot) { f.AssumptionsRef = "" }, "mixed currency": func(f *FIRESnapshot) { f.PensionAssets.Currency = "USD" },
		"negative debt": func(f *FIRESnapshot) { f.Liabilities.MinorUnits = -1 }, "target currency": func(f *FIRESnapshot) { v := Amount{1, "EUR"}; f.FIRENumber = &v },
		"zero target": func(f *FIRESnapshot) { v := Amount{0, GBP}; f.FIRENumber = &v }, "progress without target": func(f *FIRESnapshot) { f.FIRENumber = nil },
		"negative progress": func(f *FIRESnapshot) { v := int64(-1); f.ProgressBasisPoints = &v },
	} {
		t.Run(name, func(t *testing.T) { next := f; mutate(&next); mustInvalid(t, next.Validate(household, testMembers())) })
	}
	members := testMembers()
	members[2].HouseholdID = 8
	mustInvalid(t, f.Validate(household, members))
}
