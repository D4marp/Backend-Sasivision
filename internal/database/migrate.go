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

// SchemaReady reports whether core application tables exist.
func SchemaReady(db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM information_schema.tables
		WHERE table_schema = DATABASE() AND table_name = 'quiz_categories'`).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// EnsureSchema runs migrations when enabled or when the schema is missing.
func EnsureSchema(db *sql.DB, runMigrations bool) error {
	dir := MigrationsDir()
	ready, err := SchemaReady(db)
	if err != nil {
		log.Printf("[migrate] schema check failed: %v — running migrations", err)
		return RunMigrations(db, dir)
	}
	if runMigrations || !ready {
		if !ready {
			log.Println("[migrate] schema missing — running migrations")
		} else {
			log.Println("[migrate] RUN_MIGRATIONS=true — running migrations")
		}
		return RunMigrations(db, dir)
	}
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
