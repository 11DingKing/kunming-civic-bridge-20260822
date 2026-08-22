package repository

import (
	"context"
	"database/sql"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MigrationStore struct{ DB *sql.DB }

func (s MigrationStore) Applied(ctx context.Context) ([]int, error) {
	rows, e := s.DB.QueryContext(ctx, `SELECT version FROM schema_migrations ORDER BY version`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	var out []int
	for rows.Next() {
		var v int
		if e = rows.Scan(&v); e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func (s MigrationStore) Pending(files []string, applied []int) []string {
	seen := map[int]bool{}
	for _, v := range applied {
		seen[v] = true
	}
	sort.Strings(files)
	var out []string
	for _, file := range files {
		parts := strings.SplitN(file, "_", 2)
		v, e := strconv.Atoi(parts[0])
		if e == nil && !seen[v] {
			out = append(out, file)
		}
	}
	return out
}
func (s MigrationStore) Record(ctx context.Context, version int) error {
	_, e := s.DB.ExecContext(ctx, `INSERT INTO schema_migrations(version,applied_at) VALUES(?,?)`, version, time.Now().UTC().Format(time.RFC3339Nano))
	return e
}
