package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all .sql files in dir (sorted by filename).
func RunMigrations(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %q: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)

	if len(files) == 0 {
		return fmt.Errorf("no migration files in %q", dir)
	}

	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}

		sqlText := strings.TrimSpace(string(body))
		if sqlText == "" {
			continue
		}

		log.Printf("[migrate] applying %s", filepath.Base(file))
		if _, err := db.Exec(sqlText); err != nil {
			return fmt.Errorf("apply %s: %w", filepath.Base(file), err)
		}
	}

	log.Printf("[migrate] done (%d files)", len(files))
	return nil
}

// MigrationsDir resolves the migrations folder for local dev and Docker.
func MigrationsDir() string {
	if v := os.Getenv("MIGRATIONS_DIR"); v != "" {
		return v
	}
	candidates := []string{"migrations", "/app/migrations", "./migrations"}
	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}
	return "migrations"
}
