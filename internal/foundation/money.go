package foundation

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Money represents a monetary value using fixed-precision integer arithmetic.
// The amount is stored as the smallest currency unit (e.g., cents for USD,
// cents for SGD). No binary floating-point is authoritative for payroll values.
//
// Requirements: 5.2, 21.10–21.11, 23.9–23.13
type Money struct {
	// Amount in minor units (e.g., cents). May be negative for deductions.
	Amount int64

	// Currency is the ISO 4217 three-letter code.
	Currency CurrencyCode
}

// NewMoney creates a Money value with validation.
func NewMoney(amount int64, currency CurrencyCode) (Money, error) {
	if err := currency.Validate(); err != nil {
		return Money{}, err
	}
	return Money{Amount: amount, Currency: currency}, nil
}

// Add returns the sum of two Money values. Returns an error if currencies
// do not match or if overflow would occur.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.Amount + other.Amount
	// Check for overflow.
	if (other.Amount > 0 && result < m.Amount) || (other.Amount < 0 && result > m.Amount) {
		return Money{}, errors.New("monetary arithmetic overflow")
	}
	return Money{Amount: result, Currency: m.Currency}, nil
}

// Subtract returns the difference of two Money values.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, ErrCurrencyMismatch
	}
	result := m.Amount - other.Amount
	if (other.Amount > 0 && result > m.Amount) || (other.Amount < 0 && result < m.Amount) {
		return Money{}, errors.New("monetary arithmetic overflow")
	}
	return Money{Amount: result, Currency: m.Currency}, nil
}

// Multiply scales the amount by a fixed-precision multiplier expressed as
// numerator/denominator to avoid floating-point. Result is rounded according
// to the specified mode.
func (m Money) Multiply(numerator, denominator int64, mode RoundingMode) (Money, error) {
	if denominator == 0 {
		return Money{}, errors.New("division by zero")
	}
	// Use 128-bit-style calculation via intermediate.
	product := m.Amount * numerator
	// Check for overflow in the multiplication.
	if numerator != 0 && product/numerator != m.Amount {
		return Money{}, errors.New("monetary arithmetic overflow")
	}
	result := roundDivide(product, denominator, mode)
	return Money{Amount: result, Currency: m.Currency}, nil
}

// IsZero returns true for a zero amount.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsNegative returns true for a negative amount.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// String returns a human-readable representation like "SGD 1234" (in minor units).
func (m Money) String() string {
	return fmt.Sprintf("%s %d", m.Currency, m.Amount)
}

// ErrCurrencyMismatch indicates an operation was attempted across different currencies.
var ErrCurrencyMismatch = errors.New("currency mismatch")

// CurrencyCode is an ISO 4217 three-letter currency code.
type CurrencyCode string

const (
	CurrencySGD CurrencyCode = "SGD"
	CurrencyUSD CurrencyCode = "USD"
	CurrencyMYR CurrencyCode = "MYR"
)

// Validate ensures the currency code is a valid three-letter uppercase string.
func (c CurrencyCode) Validate() error {
	s := string(c)
	if len(s) != 3 || s != strings.ToUpper(s) {
		return errors.New("currency code must be exactly three uppercase letters")
	}
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return errors.New("currency code must be exactly three uppercase letters")
		}
	}
	return nil
}

// RoundingMode controls how fractional minor units are resolved.
type RoundingMode int

const (
	// RoundHalfUp rounds 0.5 towards positive infinity (banker's standard).
	RoundHalfUp RoundingMode = iota
	// RoundDown truncates towards zero.
	RoundDown
	// RoundUp rounds away from zero.
	RoundUp
)

// roundDivide performs integer division with the specified rounding mode.
func roundDivide(numerator, denominator int64, mode RoundingMode) int64 {
	if denominator == 0 {
		return 0
	}
	quotient := numerator / denominator
	remainder := numerator % denominator

	switch mode {
	case RoundHalfUp:
		absRemainder := remainder
		if absRemainder < 0 {
			absRemainder = -absRemainder
		}
		absDenominator := denominator
		if absDenominator < 0 {
			absDenominator = -absDenominator
		}
		if 2*absRemainder >= absDenominator {
			if (numerator >= 0) == (denominator > 0) {
				quotient++
			} else {
				quotient--
			}
		}
	case RoundDown:
		// Truncation is the default integer division behavior.
	case RoundUp:
		if remainder != 0 {
			if (numerator >= 0) == (denominator > 0) {
				quotient++
			} else {
				quotient--
			}
		}
	}
	return quotient
}

// MoneyFromMajor creates a Money value from major units (e.g., dollars)
// given the minor unit factor (e.g., 100 for cents).
func MoneyFromMajor(majorUnits int64, minorFactor int64, currency CurrencyCode) (Money, error) {
	if minorFactor <= 0 {
		return Money{}, errors.New("minor unit factor must be positive")
	}
	if err := currency.Validate(); err != nil {
		return Money{}, err
	}
	// Check for overflow.
	if majorUnits > 0 && majorUnits > math.MaxInt64/minorFactor {
		return Money{}, errors.New("monetary arithmetic overflow")
	}
	if majorUnits < 0 && majorUnits < math.MinInt64/minorFactor {
		return Money{}, errors.New("monetary arithmetic overflow")
	}
	return Money{Amount: majorUnits * minorFactor, Currency: currency}, nil
}
