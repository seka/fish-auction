package migrations

import "embed"

// FS represents a fs in the system.
//
//go:embed *.sql
var FS embed.FS
