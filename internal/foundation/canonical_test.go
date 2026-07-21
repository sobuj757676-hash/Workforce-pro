package foundation

import (
	"testing"
)

func TestCanonicalizeDeterministicFieldOrdering(t *testing.T) {
	// Two structurally identical objects with different source ordering
	// must produce the same canonical byte sequence.
	type sample struct {
		Zebra  string `json:"zebra"`
		Alpha  int    `json:"alpha"`
		Middle bool   `json:"middle"`
	}

	obj := sample{Zebra: "last", Alpha: 1, Middle: true}
	c1, err := Canonicalize(obj)
	if err != nil {
		t.Fatalf("Canonicalize() error = %v", err)
	}
	c2, err := Canonicalize(obj)
	if err != nil {
		t.Fatalf("Canonicalize() error = %v", err)
	}
	if string(c1) != string(c2) {
		t.Fatalf("repeated canonicalization differs: %q vs %q", c1, c2)
	}

	// Verify alphabetical key ordering.
	want := `{"alpha":1,"middle":true,"zebra":"last"}`
	if string(c1) != want {
		t.Fatalf("canonical = %q, want %q", c1, want)
	}
}

func TestCanonicalizeNestedObjects(t *testing.T) {
	type inner struct {
		Y int    `json:"y"`
		X string `json:"x"`
	}
	type outer struct {
		B inner `json:"b"`
		A int   `json:"a"`
	}
	obj := outer{B: inner{Y: 2, X: "hello"}, A: 1}
	canonical, err := Canonicalize(obj)
	if err != nil {
		t.Fatalf("Canonicalize() error = %v", err)
	}
	want := `{"a":1,"b":{"x":"hello","y":2}}`
	if string(canonical) != want {
		t.Fatalf("canonical = %q, want %q", canonical, want)
	}
}

func TestCanonicalizeArrayPreservesOrder(t *testing.T) {
	items := []map[string]int{{"b": 2, "a": 1}, {"d": 4, "c": 3}}
	canonical, err := Canonicalize(items)
	if err != nil {
		t.Fatalf("Canonicalize() error = %v", err)
	}
	want := `[{"a":1,"b":2},{"c":3,"d":4}]`
	if string(canonical) != want {
		t.Fatalf("canonical = %q, want %q", canonical, want)
	}
}

func TestCanonicalDigestDeterministic(t *testing.T) {
	type payload struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}
	obj := payload{Name: "test", Value: 42}
	d1, err := CanonicalDigest(obj)
	if err != nil {
		t.Fatalf("CanonicalDigest() error = %v", err)
	}
	d2, err := CanonicalDigest(obj)
	if err != nil {
		t.Fatalf("CanonicalDigest() error = %v", err)
	}
	if !d1.Equal(d2) {
		t.Fatal("same input produced different digests")
	}
	if d1.IsZero() {
		t.Fatal("digest should not be zero")
	}
}

func TestCanonicalizeRejectsNil(t *testing.T) {
	if _, err := Canonicalize(nil); err == nil {
		t.Fatal("nil input should be rejected")
	}
}
