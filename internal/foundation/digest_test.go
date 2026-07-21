package foundation

import "testing"

func TestDigestRoundTrip(t *testing.T) {
	content := []byte("hello, workforce payroll")
	d := ComputeDigest(content)
	if d.IsZero() {
		t.Fatal("digest should not be zero for non-empty content")
	}

	hex := d.String()
	if len(hex) != 64 {
		t.Fatalf("hex length = %d, want 64", len(hex))
	}

	parsed, err := ParseDigest(hex)
	if err != nil {
		t.Fatalf("ParseDigest() error = %v", err)
	}
	if !parsed.Equal(d) {
		t.Fatal("parsed digest does not equal original")
	}
}

func TestDigestDeterministic(t *testing.T) {
	content := []byte("same content produces same digest")
	d1 := ComputeDigest(content)
	d2 := ComputeDigest(content)
	if !d1.Equal(d2) {
		t.Fatal("same content produced different digests")
	}
}

func TestDigestDifferentForDifferentContent(t *testing.T) {
	d1 := ComputeDigest([]byte("content A"))
	d2 := ComputeDigest([]byte("content B"))
	if d1.Equal(d2) {
		t.Fatal("different content produced same digest")
	}
}

func TestParseDigestRejectsInvalid(t *testing.T) {
	invalid := []string{
		"",
		"not-a-hex",
		"abc123",
		"gg" + "00000000000000000000000000000000000000000000000000000000000000", // invalid hex char
	}
	for _, v := range invalid {
		if _, err := ParseDigest(v); err == nil {
			t.Fatalf("ParseDigest(%q) unexpectedly succeeded", v)
		}
	}
}
