CREATE TABLE principal_bindings (
    principal_id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    issuer TEXT NOT NULL,
    subject TEXT NOT NULL,
    principal_kind TEXT NOT NULL CHECK (principal_kind IN ('HUMAN', 'INTEGRATION')),
    status TEXT NOT NULL CHECK (status IN ('ACTIVE', 'DISABLED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT principal_bindings_authenticated_identity_unique
        UNIQUE (issuer, subject, principal_kind),
    CONSTRAINT principal_bindings_nonempty_issuer CHECK (btrim(issuer) <> ''),
    CONSTRAINT principal_bindings_nonempty_subject CHECK (btrim(subject) <> '')
);

CREATE INDEX principal_bindings_active_lookup_idx
    ON principal_bindings (issuer, subject, principal_kind)
    WHERE status = 'ACTIVE';

COMMENT ON TABLE principal_bindings IS
    'Server-authoritative mapping from a verified external identity to exactly one effective Tenant.';
COMMENT ON COLUMN principal_bindings.tenant_id IS
    'Never selected from caller-supplied request data or an unverified token claim.';
