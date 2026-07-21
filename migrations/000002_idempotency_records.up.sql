CREATE TABLE idempotency_records (
    tenant_id UUID NOT NULL,
    idempotency_key TEXT NOT NULL,
    request_digest TEXT NOT NULL,
    result_ref TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (tenant_id, idempotency_key),
    CONSTRAINT idempotency_records_nonempty_key CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT idempotency_records_nonempty_digest CHECK (btrim(request_digest) <> ''),
    CONSTRAINT idempotency_records_valid_expiry CHECK (expires_at > created_at)
);

CREATE INDEX idempotency_records_expiry_idx
    ON idempotency_records (expires_at)
    WHERE expires_at < CURRENT_TIMESTAMP;

COMMENT ON TABLE idempotency_records IS
    'Tenant-scoped idempotency store ensuring safe mutation retries converge on one effect. Requirements 7.1–7.4.';
COMMENT ON COLUMN idempotency_records.tenant_id IS
    'Every idempotency key is scoped to a Tenant; cross-Tenant collisions are impossible.';
COMMENT ON COLUMN idempotency_records.request_digest IS
    'SHA-256 hex of the canonical request payload; reuse with a different digest is rejected.';
