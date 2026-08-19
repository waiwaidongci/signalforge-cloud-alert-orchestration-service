package migrate

import (
	"database/sql"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"github.com/acme/signalforge/migrations"
)

func Apply(raw *sql.DB, driver string) error {
	folder := "sqlite"
	if driver == "postgres" {
		folder = "postgres"
	}
	entries, err := fs.ReadDir(migrations.FS, folder)
	if err != nil {
		return fmt.Errorf("read migrations %s: %w", folder, err)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		rawSQL, err := migrations.FS.ReadFile(path.Join(folder, entry.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", entry.Name(), err)
		}
		if _, err := raw.Exec(string(rawSQL)); err != nil {
			return fmt.Errorf("apply migration %s: %w", entry.Name(), err)
		}
	}
	return nil
}
