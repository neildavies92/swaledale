package domain

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Money stores GBP values as integer pence to keep totals deterministic.
type Money int64

func Pounds(amount float64) Money {
	return Money(math.Round(amount * 100))
}

func (m Money) Float64() float64 {
	return float64(m) / 100
}

func (m Money) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(m.Float64(), 'f', 2, 64)), nil
}

func (m *Money) UnmarshalJSON(data []byte) error {
	var number float64
	if err := json.Unmarshal(data, &number); err != nil {
		return err
	}
	*m = Pounds(number)
	return nil
}

func (m Money) String() string {
	sign := ""
	value := int64(m)
	if value < 0 {
		sign = "-"
		value = -value
	}
	return fmt.Sprintf("%s£%d.%02d", sign, value/100, value%100)
}
