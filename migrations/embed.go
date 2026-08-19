package migrations

import "embed"

// FS contains SQL migrations for both supported database engines.
//
//go:embed sqlite/*.sql postgres/*.sql
var FS embed.FS
