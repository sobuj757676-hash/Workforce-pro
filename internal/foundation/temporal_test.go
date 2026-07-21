package foundation

import (
	"testing"
	"time"
)

func TestSiteTimePreservesZoneAndOffset(t *testing.T) {
	instant := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	st, err := NewSiteTime(instant, "Asia/Singapore")
	if err != nil {
		t.Fatalf("NewSiteTime() error = %v", err)
	}
	if st.Zone != "Asia/Singapore" {
		t.Fatalf("Zone = %q, want Asia/Singapore", st.Zone)
	}
	// Singapore is UTC+8.
	if st.Offset != 8*3600 {
		t.Fatalf("Offset = %d, want %d", st.Offset, 8*3600)
	}
	local := st.LocalTime()
	if local.Hour() != 18 {
		t.Fatalf("LocalTime().Hour() = %d, want 18", local.Hour())
	}
}

func TestSiteTimeRejectsInvalidZone(t *testing.T) {
	instant := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	if _, err := NewSiteTime(instant, "Invalid/Zone"); err == nil {
		t.Fatal("accepted invalid time zone")
	}
}

func TestSiteTimeRejectsZeroInstant(t *testing.T) {
	if _, err := NewSiteTime(time.Time{}, "Asia/Singapore"); err == nil {
		t.Fatal("accepted zero instant")
	}
}

func TestEffectiveIntervalContains(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	ei, err := NewEffectiveInterval(start, &end)
	if err != nil {
		t.Fatalf("NewEffectiveInterval() error = %v", err)
	}

	mid := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	if !ei.Contains(mid) {
		t.Fatal("interval should contain mid-point")
	}
	if !ei.Contains(start) {
		t.Fatal("interval should contain start (inclusive)")
	}
	if ei.Contains(end) {
		t.Fatal("interval should not contain end (exclusive)")
	}
	before := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	if ei.Contains(before) {
		t.Fatal("interval should not contain before-start")
	}
}

func TestEffectiveIntervalOpenEnded(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ei, err := NewEffectiveInterval(start, nil)
	if err != nil {
		t.Fatalf("NewEffectiveInterval() error = %v", err)
	}
	if !ei.IsOpenEnded() {
		t.Fatal("should be open-ended")
	}
	future := time.Date(2099, 12, 31, 0, 0, 0, 0, time.UTC)
	if !ei.Contains(future) {
		t.Fatal("open-ended interval should contain any future date")
	}
}

func TestEffectiveIntervalOverlaps(t *testing.T) {
	start1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	end1 := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	ei1, _ := NewEffectiveInterval(start1, &end1)

	// Overlapping interval.
	start2 := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	end2 := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	ei2, _ := NewEffectiveInterval(start2, &end2)
	if !ei1.Overlaps(ei2) {
		t.Fatal("overlapping intervals should report overlap")
	}

	// Non-overlapping interval (exactly adjacent).
	start3 := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	end3 := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	ei3, _ := NewEffectiveInterval(start3, &end3)
	if ei1.Overlaps(ei3) {
		t.Fatal("adjacent intervals should not overlap (half-open)")
	}

	// Open-ended overlaps with anything after its start.
	eiOpen, _ := NewEffectiveInterval(start1, nil)
	if !eiOpen.Overlaps(ei2) {
		t.Fatal("open-ended interval should overlap with any starting after its start")
	}
}

func TestEffectiveIntervalValidation(t *testing.T) {
	// Zero start.
	if _, err := NewEffectiveInterval(time.Time{}, nil); err == nil {
		t.Fatal("accepted zero start")
	}
	// End before start.
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := NewEffectiveInterval(start, &end); err == nil {
		t.Fatal("accepted end before start")
	}
}

func TestDurationValidation(t *testing.T) {
	d, err := NewDuration(480)
	if err != nil {
		t.Fatalf("NewDuration() error = %v", err)
	}
	if d.Minutes != 480 || d.IsZero() {
		t.Fatalf("unexpected duration: %+v", d)
	}

	zero, _ := NewDuration(0)
	if !zero.IsZero() {
		t.Fatal("zero duration should be zero")
	}

	if _, err := NewDuration(-1); err == nil {
		t.Fatal("accepted negative duration")
	}
}

func TestDurationAdd(t *testing.T) {
	d1, _ := NewDuration(60)
	d2, _ := NewDuration(30)
	sum := d1.Add(d2)
	if sum.Minutes != 90 {
		t.Fatalf("Add() = %d, want 90", sum.Minutes)
	}
}

func TestVersionedReferenceValidation(t *testing.T) {
	tenantID := MustParseIdentifier("018f6f4a-7b2c-7d11-8a33-123456789abc")
	entityID := MustParseIdentifier("018f6f4a-7b2c-7d11-bb66-556677889900")
	scopedID, _ := NewTenantScopedID(tenantID, entityID)

	ref := VersionedReference{
		AggregateType: "CompensationProfile",
		ID:            scopedID,
		Version:       3,
		ContentDigest: ComputeDigest([]byte("content")),
	}
	if err := ref.Validate(); err != nil {
		t.Fatalf("valid reference rejected: %v", err)
	}

	// Missing aggregate type.
	ref2 := ref
	ref2.AggregateType = ""
	if err := ref2.Validate(); err == nil {
		t.Fatal("accepted empty aggregate type")
	}

	// Zero version.
	ref3 := ref
	ref3.Version = 0
	if err := ref3.Validate(); err == nil {
		t.Fatal("accepted zero version")
	}

	// Zero digest.
	ref4 := ref
	ref4.ContentDigest = Digest{}
	if err := ref4.Validate(); err == nil {
		t.Fatal("accepted zero digest")
	}
}
