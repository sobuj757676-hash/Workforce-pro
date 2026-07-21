package foundation

import (
	"encoding/json"
	"errors"
	"sort"
)

// ErrCanonicalizeInput is returned when input cannot be canonicalized.
var ErrCanonicalizeInput = errors.New("cannot canonicalize input")

// Canonicalize produces a deterministic byte representation of a value
// suitable for digest computation. The output uses sorted JSON keys and
// fixed representations so that semantically identical content always
// produces the same byte sequence.
//
// Requirements: 6.5–6.6, 7.8–7.9, 13.3, 14.2, 18.2, 19.3, 23.13, 25.10–25.14
func Canonicalize(value any) ([]byte, error) {
	if value == nil {
		return nil, ErrCanonicalizeInput
	}

	// Convert to a JSON-compatible intermediate to normalize all types.
	intermediate, err := json.Marshal(value)
	if err != nil {
		return nil, ErrCanonicalizeInput
	}

	// Parse into the generic interface to sort keys deterministically.
	var parsed any
	if err := json.Unmarshal(intermediate, &parsed); err != nil {
		return nil, ErrCanonicalizeInput
	}

	canonical := canonicalizeValue(parsed)
	result, err := json.Marshal(canonical)
	if err != nil {
		return nil, ErrCanonicalizeInput
	}
	return result, nil
}

// CanonicalDigest computes the SHA-256 digest of the canonical
// representation of a value.
func CanonicalDigest(value any) (Digest, error) {
	canonical, err := Canonicalize(value)
	if err != nil {
		return Digest{}, err
	}
	return ComputeDigest(canonical), nil
}

// canonicalizeValue recursively sorts object keys and normalizes values.
func canonicalizeValue(value any) any {
	switch v := value.(type) {
	case map[string]any:
		return canonicalizeObject(v)
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = canonicalizeValue(item)
		}
		return result
	default:
		return v
	}
}

// orderedEntry preserves insertion order for JSON serialization.
type orderedEntry struct {
	Key   string
	Value any
}

// orderedMap serializes as a JSON object with deterministic field ordering.
type orderedMap struct {
	Entries []orderedEntry
}

func (om orderedMap) MarshalJSON() ([]byte, error) {
	buf := []byte("{")
	for i, entry := range om.Entries {
		if i > 0 {
			buf = append(buf, ',')
		}
		key, err := json.Marshal(entry.Key)
		if err != nil {
			return nil, err
		}
		val, err := json.Marshal(entry.Value)
		if err != nil {
			return nil, err
		}
		buf = append(buf, key...)
		buf = append(buf, ':')
		buf = append(buf, val...)
	}
	buf = append(buf, '}')
	return buf, nil
}

func canonicalizeObject(obj map[string]any) orderedMap {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]orderedEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, orderedEntry{
			Key:   k,
			Value: canonicalizeValue(obj[k]),
		})
	}
	return orderedMap{Entries: entries}
}
