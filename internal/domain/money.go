package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

// Money is a signed integer count of minor currency units. Canonical values pair
// it with Currency in Amount. Bare Money's JSON and String methods are retained
// only for the legacy GBP budget API (decimal pounds, two decimal places).
type Money int64

func (m Money) decimalPounds() string {
	n := int64(m)
	sign := ""
	magnitude := uint64(n)
	if n < 0 {
		sign = "-"
		magnitude = uint64(-(n + 1)) + 1
	}
	return fmt.Sprintf("%s%d.%02d", sign, magnitude/100, magnitude%100)
}

func (m Money) MarshalJSON() ([]byte, error) { return []byte(m.decimalPounds()), nil }

// UnmarshalJSON parses legacy decimal pounds exactly, rejecting sub-penny values
// and overflow instead of rounding through float64.
func (m *Money) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if !json.Valid(data) || len(data) == 0 || !(data[0] == '-' || data[0] >= '0' && data[0] <= '9') {
		return fmt.Errorf("money must be a JSON number")
	}
	value, ok := new(big.Rat).SetString(string(data))
	if !ok {
		return fmt.Errorf("invalid money")
	}
	value.Mul(value, big.NewRat(100, 1))
	if !value.IsInt() || !value.Num().IsInt64() {
		return fmt.Errorf("money must fit int64 minor units without rounding")
	}
	*m = Money(value.Num().Int64())
	return nil
}

func (m Money) String() string {
	value := m.decimalPounds()
	if m < 0 {
		return "-£" + value[1:]
	}
	return "£" + value
}

type Currency string

const GBP Currency = "GBP"

// Validate checks the ISO-style code shape. Supported currencies and their minor
// unit exponents are connector/configuration concerns, not an FX engine here.
func (c Currency) Validate() error {
	if len(c) != 3 {
		return fmt.Errorf("currency must have three uppercase letters")
	}
	for _, ch := range c {
		if ch < 'A' || ch > 'Z' {
			return fmt.Errorf("invalid currency %q", c)
		}
	}
	return nil
}

// Amount is currency-qualified Money. Its wire form is integer minor units,
// independent of legacy Money's decimal-GBP JSON convention.
type Amount struct {
	MinorUnits Money    `json:"-"`
	Currency   Currency `json:"currency"`
}

func (a Amount) Validate() error { return a.Currency.Validate() }
func (a Amount) MarshalJSON() ([]byte, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}
	return []byte(`{"minorUnits":` + strconv.FormatInt(int64(a.MinorUnits), 10) + `,"currency":"` + string(a.Currency) + `"}`), nil
}
func (a *Amount) UnmarshalJSON(data []byte) error {
	var wire struct {
		MinorUnits *int64   `json:"minorUnits"`
		Currency   Currency `json:"currency"`
	}
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	if wire.MinorUnits == nil {
		return fmt.Errorf("minorUnits is required")
	}
	next := Amount{Money(*wire.MinorUnits), wire.Currency}
	if err := next.Validate(); err != nil {
		return err
	}
	*a = next
	return nil
}
