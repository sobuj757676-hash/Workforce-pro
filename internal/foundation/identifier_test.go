package foundation

import "testing"

func TestIdentifierRoundTrip(t *testing.T) {
	const value = "018f6f4a-7b2c-7d11-8a33-123456789abc"
	id, err := ParseIdentifier(value)
	if err != nil {
		t.Fatalf("ParseIdentifier() error = %v", err)
	}
	if got := id.String(); got != value {
		t.Fatalf("Identifier.String() = %q, want %q", got, value)
	}
}

func TestIdentifierRejectsInvalidValues(t *testing.T) {
	for _, value := range []string{"", "not-a-uuid", "00000000-0000-0000-0000-000000000000", "018f6f4a-7b2c-7d11-0a33-123456789abc"} {
		if _, err := ParseIdentifier(value); err == nil {
			t.Fatalf("ParseIdentifier(%q) unexpectedly succeeded", value)
		}
	}
}
