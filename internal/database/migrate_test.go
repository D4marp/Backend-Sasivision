package database

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrationsDir(t *testing.T) {
	dir := MigrationsDir()
	if dir == "" {
		t.Fatal("expected non-empty migrations dir")
	}
	if _, err := os.Stat(dir); err != nil {
		// When running from package dir, try repo root relative path
		repo := filepath.Join("..", "..", "migrations")
		if _, err2 := os.Stat(repo); err2 != nil {
			t.Skipf("migrations folder not found: %v", err)
		}
	}
}

func TestRunMigrationsRequiresDir(t *testing.T) {
	err := RunMigrations(nil, filepath.Join(t.TempDir(), "empty"))
	if err == nil {
		t.Fatal("expected error for empty migrations dir")
	}
}
