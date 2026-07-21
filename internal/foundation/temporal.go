package foundation

import (
	"errors"
	"fmt"
	"time"
)

// SiteTime represents a moment in time with an explicit IANA time zone
// and resolved UTC offset, preserving the operational context of the Site.
//
// Requirements: 5.2, 17.12–17.13, 21.5–21.11
type SiteTime struct {
	// Instant is the absolute UTC moment.
	Instant time.Time

	// Zone is the IANA time zone name (e.g., "Asia/Singapore").
	Zone string

	// Offset is the resolved UTC offset in seconds at the instant.
	Offset int
}

// NewSiteTime creates a SiteTime with validation.
func NewSiteTime(instant time.Time, zone string) (SiteTime, error) {
	if instant.IsZero() {
		return SiteTime{}, errors.New("instant is required")
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return SiteTime{}, fmt.Errorf("invalid IANA time zone %q: %w", zone, err)
	}
	_, offset := instant.In(loc).Zone()
	return SiteTime{
		Instant: instant.UTC(),
		Zone:    zone,
		Offset:  offset,
	}, nil
}

// LocalTime returns the time expressed in the Site's time zone.
func (st SiteTime) LocalTime() time.Time {
	loc, err := time.LoadLocation(st.Zone)
	if err != nil {
		return st.Instant
	}
	return st.Instant.In(loc)
}

// EffectiveInterval represents a half-open time range [starts_at, ends_at)
// for effective-dated records like compensation profiles and assignments.
// A nil ends_at represents "currently effective" (no known end).
//
// Requirements: 21.5–21.11
type EffectiveInterval struct {
	StartsAt time.Time
	EndsAt   *time.Time
}

// NewEffectiveInterval creates and validates an interval.
func NewEffectiveInterval(startsAt time.Time, endsAt *time.Time) (EffectiveInterval, error) {
	if startsAt.IsZero() {
		return EffectiveInterval{}, errors.New("starts_at is required")
	}
	if endsAt != nil && !endsAt.After(startsAt) {
		return EffectiveInterval{}, errors.New("ends_at must be after starts_at")
	}
	return EffectiveInterval{StartsAt: startsAt, EndsAt: endsAt}, nil
}

// Contains returns true if the given instant falls within the interval.
func (ei EffectiveInterval) Contains(instant time.Time) bool {
	if instant.Before(ei.StartsAt) {
		return false
	}
	if ei.EndsAt != nil && !instant.Before(*ei.EndsAt) {
		return false
	}
	return true
}

// Overlaps returns true if this interval overlaps with another.
func (ei EffectiveInterval) Overlaps(other EffectiveInterval) bool {
	// Two half-open intervals [a, b) and [c, d) overlap iff a < d and c < b.
	if ei.EndsAt != nil && !other.StartsAt.Before(*ei.EndsAt) {
		return false
	}
	if other.EndsAt != nil && !ei.StartsAt.Before(*other.EndsAt) {
		return false
	}
	return true
}

// IsOpenEnded returns true if the interval has no defined end.
func (ei EffectiveInterval) IsOpenEnded() bool {
	return ei.EndsAt == nil
}

// Duration is a non-negative work duration in minutes.
// Used for break times, payable durations, and planned work windows.
type Duration struct {
	Minutes int64
}

// NewDuration creates a validated non-negative duration.
func NewDuration(minutes int64) (Duration, error) {
	if minutes < 0 {
		return Duration{}, errors.New("duration cannot be negative")
	}
	return Duration{Minutes: minutes}, nil
}

// IsZero returns true for a zero-length duration.
func (d Duration) IsZero() bool {
	return d.Minutes == 0
}

// Add combines two durations.
func (d Duration) Add(other Duration) Duration {
	return Duration{Minutes: d.Minutes + other.Minutes}
}

// VersionedReference is a typed reference to an exact version of another
// domain aggregate. Used to lock provenance in calculations and snapshots.
//
// Requirements: 23.13, 25.10–25.14
type VersionedReference struct {
	// AggregateType identifies the kind of referenced entity
	// (e.g., "CalendarVersion", "CompensationProfile", "PayRule").
	AggregateType string

	// ID is the Tenant-scoped identifier of the referenced entity.
	ID TenantScopedID

	// Version is the exact version number referenced.
	Version int64

	// ContentDigest is the digest of the referenced version's content.
	ContentDigest Digest
}

// Validate ensures the reference is complete.
func (vr VersionedReference) Validate() error {
	if vr.AggregateType == "" {
		return errors.New("aggregate_type is required")
	}
	if err := vr.ID.Validate(); err != nil {
		return fmt.Errorf("referenced id: %w", err)
	}
	if vr.Version <= 0 {
		return errors.New("version must be positive")
	}
	if vr.ContentDigest.IsZero() {
		return errors.New("content_digest is required")
	}
	return nil
}
