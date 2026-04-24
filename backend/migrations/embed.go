package migrations

import "embed"

// FS is used by golang-migrate. Contains *.up.sql in the root directory.
// Rollback scripts are in the rollback/ subdirectory for manual use.
//
//go:embed *.sql
var FS embed.FS
