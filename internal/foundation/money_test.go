package foundation

import (
	"testing"
)

func TestMoneyCreationAndValidation(t *testing.T) {
	m, err := NewMoney(15000, CurrencySGD)
	if err != nil {
		t.Fatalf("NewMoney() error = %v", err)
	}
	if m.Amount != 15000 || m.Currency != CurrencySGD {
		t.Fatalf("unexpected Money: %+v", m)
	}
	if m.IsZero() || m.IsNegative() {
		t.Fatal("unexpected zero/negative state")
	}
}

func TestMoneyRejectsInvalidCurrency(t *testing.T) {
	if _, err := NewMoney(100, "sg"); err == nil {
		t.Fatal("accepted invalid currency")
	}
	if _, err := NewMoney(100, "SGDD"); err == nil {
		t.Fatal("accepted 4-letter currency")
	}
	if _, err := NewMoney(100, "123"); err == nil {
		t.Fatal("accepted numeric currency")
	}
}

func TestMoneyAddition(t *testing.T) {
	a, _ := NewMoney(10000, CurrencySGD)
	b, _ := NewMoney(5000, CurrencySGD)
	result, err := a.Add(b)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}
	if result.Amount != 15000 {
		t.Fatalf("Add() = %d, want 15000", result.Amount)
	}
}

func TestMoneyAdditionCurrencyMismatch(t *testing.T) {
	a, _ := NewMoney(10000, CurrencySGD)
	b, _ := NewMoney(5000, CurrencyUSD)
	if _, err := a.Add(b); err != ErrCurrencyMismatch {
		t.Fatalf("expected currency mismatch, got %v", err)
	}
}

func TestMoneySubtraction(t *testing.T) {
	a, _ := NewMoney(10000, CurrencySGD)
	b, _ := NewMoney(3000, CurrencySGD)
	result, err := a.Subtract(b)
	if err != nil {
		t.Fatalf("Subtract() error = %v", err)
	}
	if result.Amount != 7000 {
		t.Fatalf("Subtract() = %d, want 7000", result.Amount)
	}
}

func TestMoneyMultiply(t *testing.T) {
	m, _ := NewMoney(10000, CurrencySGD) // $100.00

	// 1.5x with half-up rounding.
	result, err := m.Multiply(3, 2, RoundHalfUp)
	if err != nil {
		t.Fatalf("Multiply() error = %v", err)
	}
	if result.Amount != 15000 {
		t.Fatalf("Multiply(3/2) = %d, want 15000", result.Amount)
	}

	// Rounding test: 10001 * 1/3 = 3333.67 -> 3334 half-up.
	m2, _ := NewMoney(10001, CurrencySGD)
	result2, err := m2.Multiply(1, 3, RoundHalfUp)
	if err != nil {
		t.Fatalf("Multiply() error = %v", err)
	}
	if result2.Amount != 3334 {
		t.Fatalf("Multiply(1/3) of 10001 = %d, want 3334", result2.Amount)
	}

	// RoundDown: 10001 * 1/3 = 3333.
	result3, err := m2.Multiply(1, 3, RoundDown)
	if err != nil {
		t.Fatalf("Multiply() error = %v", err)
	}
	if result3.Amount != 3333 {
		t.Fatalf("Multiply(1/3) RoundDown = %d, want 3333", result3.Amount)
	}
}

func TestMoneyFromMajor(t *testing.T) {
	m, err := MoneyFromMajor(150, 100, CurrencySGD)
	if err != nil {
		t.Fatalf("MoneyFromMajor() error = %v", err)
	}
	if m.Amount != 15000 {
		t.Fatalf("MoneyFromMajor(150, 100) = %d, want 15000", m.Amount)
	}
}

func TestCurrencyCodeValidation(t *testing.T) {
	valid := []CurrencyCode{"SGD", "USD", "MYR", "EUR", "GBP"}
	for _, c := range valid {
		if err := c.Validate(); err != nil {
			t.Fatalf("valid currency %q rejected: %v", c, err)
		}
	}
	invalid := []CurrencyCode{"", "sg", "SGDD", "12D", "S D"}
	for _, c := range invalid {
		if err := c.Validate(); err == nil {
			t.Fatalf("invalid currency %q accepted", c)
		}
	}
}
