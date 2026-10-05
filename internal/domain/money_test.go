package domain

import (
	"encoding/json"
	"math"
	"testing"
)

func TestLegacyMoneyExactRoundTrip(t *testing.T) {
	for _, m := range []Money{0, 1, -1, 12345, 9007199254740993, Money(math.MaxInt64), Money(math.MinInt64)} {
		data, err := json.Marshal(m)
		mustValid(t, err)
		var restored Money
		mustValid(t, json.Unmarshal(data, &restored))
		if restored != m {
			t.Fatalf("money precision lost: %d -> %s -> %d", int64(m), data, int64(restored))
		}
	}
	for _, bad := range []string{`null`, `"12.34"`, `true`, `{}`, `0.001`, `92233720368547758.08`, `-92233720368547758.09`} {
		value := Money(12)
		mustInvalid(t, json.Unmarshal([]byte(bad), &value))
		if value != 12 {
			t.Fatal("failed decode changed money")
		}
	}
	var exponent Money
	mustValid(t, json.Unmarshal([]byte(`1e2`), &exponent))
	if exponent != 10000 {
		t.Fatal("wrong exponent value")
	}
}
func TestAmountIntegerWireAndCurrency(t *testing.T) {
	for _, c := range []Currency{GBP, "EUR", "USD", "JPY", "KWD"} {
		a := Amount{Money(math.MaxInt64), c}
		data, err := json.Marshal(a)
		mustValid(t, err)
		var wire struct {
			MinorUnits int64 `json:"minorUnits"`
		}
		mustValid(t, json.Unmarshal(data, &wire))
		if wire.MinorUnits != math.MaxInt64 {
			t.Fatal("canonical minor units were scaled")
		}
		var restored Amount
		mustValid(t, json.Unmarshal(data, &restored))
		if a != restored {
			t.Fatal("amount roundtrip lost information")
		}
	}
	for _, bad := range []string{`{}`, `null`, `{"minorUnits":null,"currency":"GBP"}`, `{"minorUnits":1.2,"currency":"GBP"}`, `{"minorUnits":1,"currency":"gbp"}`, `{"minorUnits":9223372036854775808,"currency":"GBP"}`} {
		a := Amount{1, GBP}
		mustInvalid(t, json.Unmarshal([]byte(bad), &a))
		if a != (Amount{1, GBP}) {
			t.Fatal("failed decode changed amount")
		}
	}
	for _, c := range []Currency{"", "G", "GBP4", "gbp", "GB£"} {
		mustInvalid(t, c.Validate())
	}
	_, err := json.Marshal(Amount{})
	mustInvalid(t, err)
}
