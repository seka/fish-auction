package migrations

import "embed"

//go:embed *.sql
// FS represents a fs in the system.
var FS embed.FS
