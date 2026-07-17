package foundation

import (
	"encoding/hex"
	"errors"
	"strings"
)

var ErrInvalidIdentifier = errors.New("invalid identifier")

// Identifier is a canonical RFC 4122 UUID representation used at trust boundaries.
type Identifier [16]byte

func ParseIdentifier(value string) (Identifier, error) {
	var id Identifier
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return id, ErrInvalidIdentifier
	}

	compact := strings.ReplaceAll(value, "-", "")
	decoded, err := hex.DecodeString(compact)
	if err != nil || len(decoded) != len(id) {
		return id, ErrInvalidIdentifier
	}
	copy(id[:], decoded)

	// Only RFC 4122 variants are accepted. Nil identifiers are never authoritative IDs.
	if id == (Identifier{}) || id[8]&0xc0 != 0x80 {
		return Identifier{}, ErrInvalidIdentifier
	}
	return id, nil
}

func MustParseIdentifier(value string) Identifier {
	id, err := ParseIdentifier(value)
	if err != nil {
		panic(err)
	}
	return id
}

func (id Identifier) String() string {
	compact := hex.EncodeToString(id[:])
	return compact[0:8] + "-" + compact[8:12] + "-" + compact[12:16] + "-" + compact[16:20] + "-" + compact[20:32]
}

func (id Identifier) IsZero() bool {
	return id == (Identifier{})
}
