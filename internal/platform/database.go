package platform

import (
	"context"
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func OpenDatabase(ctx context.Context, path string) (*sql.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if e := os.MkdirAll(dir, 0755); e != nil {
			return nil, e
		}
	}
	db, e := sql.Open("sqlite", fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", path))
	if e != nil {
		return nil, e
	}
	db.SetMaxOpenConns(8)
	if e = db.PingContext(ctx); e != nil {
		db.Close()
		return nil, e
	}
	return db, nil
}
func Migrate(ctx context.Context, db *sql.DB) error {
	root := "migrations"
	if _, err := os.Stat(root); err != nil {
		root = "../../migrations"
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)
	if _, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations(version INTEGER PRIMARY KEY, applied_at TEXT NOT NULL)`); err != nil {
		return err
	}
	for _, file := range files {
		parts := strings.SplitN(file, "_", 2)
		version, parseErr := strconv.Atoi(parts[0])
		if parseErr != nil {
			return parseErr
		}
		var applied int
		if err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM schema_migrations WHERE version=?`, version).Scan(&applied); err != nil {
			return err
		}
		if applied > 0 {
			continue
		}
		content, readErr := os.ReadFile(filepath.Join(root, file))
		if readErr != nil {
			return readErr
		}
		tx, beginErr := db.BeginTx(ctx, nil)
		if beginErr != nil {
			return beginErr
		}
		if _, err = tx.ExecContext(ctx, string(content)); err != nil {
			tx.Rollback()
			return err
		}
		if _, err = tx.ExecContext(ctx, `INSERT INTO schema_migrations(version,applied_at) VALUES(?,?)`, version, time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
			tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}
