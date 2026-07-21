package foundation

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// ErrInvalidDigest is returned when a digest string is malformed.
var ErrInvalidDigest = errors.New("invalid digest")

// Digest is a SHA-256 content digest used to bind approvals, signatures,
// and version comparisons to exact canonical content.
// Requirements: 6.5–6.6, 7.8–7.9
type Digest [32]byte

// ComputeDigest produces a SHA-256 digest from canonical byte content.
func ComputeDigest(content []byte) Digest {
	return sha256.Sum256(content)
}

// ParseDigest parses a lowercase hex-encoded SHA-256 string.
func ParseDigest(value string) (Digest, error) {
	if len(value) != 64 {
		return Digest{}, ErrInvalidDigest
	}
	decoded, err := hex.DecodeString(value)
	if err != nil || len(decoded) != 32 {
		return Digest{}, ErrInvalidDigest
	}
	var d Digest
	copy(d[:], decoded)
	return d, nil
}

// String returns the lowercase hex representation.
func (d Digest) String() string {
	return hex.EncodeToString(d[:])
}

// IsZero returns true if the digest is the zero value.
func (d Digest) IsZero() bool {
	return d == (Digest{})
}

// Equal compares two digests in constant time via byte comparison.
func (d Digest) Equal(other Digest) bool {
	return d == other
}
