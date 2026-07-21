package migrations

import _ "embed"

//go:embed 000001_principal_bindings.up.sql
var PrincipalBindingsUp string

//go:embed 000001_principal_bindings.down.sql
var PrincipalBindingsDown string

//go:embed 000002_idempotency_records.up.sql
var IdempotencyRecordsUp string

//go:embed 000002_idempotency_records.down.sql
var IdempotencyRecordsDown string
