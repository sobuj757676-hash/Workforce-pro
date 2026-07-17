package migrations

import _ "embed"

//go:embed 000001_principal_bindings.up.sql
var PrincipalBindingsUp string

//go:embed 000001_principal_bindings.down.sql
var PrincipalBindingsDown string
