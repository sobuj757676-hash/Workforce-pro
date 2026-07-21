package foundation

import (
	"testing"
	"time"
)

func validCommandEnvelope() CommandEnvelope {
	return CommandEnvelope{
		PrincipalID:    MustParseIdentifier("018f6f4a-7b2c-7d11-aa55-001122334455"),
		TenantID:       MustParseIdentifier("018f6f4a-7b2c-7d11-8a33-123456789abc"),
		PrincipalKind:  "HUMAN",
		Permission:     "calendar:publish",
		CorrelationID:  "corr-001",
		IdempotencyKey: "idem-001",
		RequestDigest:  ComputeDigest([]byte("test-payload")),
		IssuedAt:       time.Now(),
	}
}

func TestCommandEnvelopeValidation(t *testing.T) {
	cmd := validCommandEnvelope()
	if err := cmd.Validate(); err != nil {
		t.Fatalf("valid envelope rejected: %v", err)
	}
}

func TestCommandEnvelopeRequiresMandatoryFields(t *testing.T) {
	tests := map[string]func(*CommandEnvelope){
		"missing principal":       func(cmd *CommandEnvelope) { cmd.PrincipalID = Identifier{} },
		"missing tenant":          func(cmd *CommandEnvelope) { cmd.TenantID = Identifier{} },
		"missing principal_kind":  func(cmd *CommandEnvelope) { cmd.PrincipalKind = "" },
		"missing correlation_id":  func(cmd *CommandEnvelope) { cmd.CorrelationID = "" },
		"missing idempotency_key": func(cmd *CommandEnvelope) { cmd.IdempotencyKey = "" },
		"missing request_digest":  func(cmd *CommandEnvelope) { cmd.RequestDigest = Digest{} },
		"missing issued_at":       func(cmd *CommandEnvelope) { cmd.IssuedAt = time.Time{} },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := validCommandEnvelope()
			mutate(&cmd)
			if err := cmd.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", name)
			}
		})
	}
}

func TestCommandEnvelopeVersionCheck(t *testing.T) {
	cmd := validCommandEnvelope()

	// No expected version: always passes.
	if err := cmd.CheckVersion(5); err != nil {
		t.Fatalf("nil expected version should pass: %v", err)
	}

	// Matching version: passes.
	v := int64(5)
	cmd.ExpectedVersion = &v
	if err := cmd.CheckVersion(5); err != nil {
		t.Fatalf("matching version should pass: %v", err)
	}

	// Mismatched version: conflict.
	if err := cmd.CheckVersion(6); err != ErrConflict {
		t.Fatalf("mismatched version should return ErrConflict, got %v", err)
	}
}
