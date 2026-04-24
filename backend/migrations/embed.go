package migrations

import "embed"

// FS is used by golang-migrate. Contains *.up.sql and *.down.sql.
//
//go:embed *.sql
var FS embed.FS
