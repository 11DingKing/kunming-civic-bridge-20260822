package repository

import (
	"context"
	"database/sql"
	"github.com/11DingKing/kunming-civic-bridge/internal/domain"
	"github.com/11DingKing/kunming-civic-bridge/internal/platform"
	"testing"
)

func identityDB(t *testing.T) *sql.DB {
	db, e := sql.Open("sqlite", ":memory:")
	if e != nil {
		t.Fatal(e)
	}
	db.SetMaxOpenConns(1)
	if e = platform.Migrate(context.Background(), db); e != nil {
		t.Fatal(e)
	}
	return db
}
func TestIdentityStoreAuthentication(t *testing.T) {
	db := identityDB(t)
	defer db.Close()
	s := IdentityStore{DB: db}
	id, e := s.Create(context.Background(), "市民", "13800000001", "secret", "云南/昆明/西山/马街", domain.RoleCitizen)
	if e != nil {
		t.Fatal(e)
	}
	u, e := s.Authenticate(context.Background(), "13800000001", "secret")
	if e != nil || u.ID != id || len(u.Roles) != 1 {
		t.Fatalf("user=%+v err=%v", u, e)
	}
	if _, e = s.Authenticate(context.Background(), "13800000001", "wrong"); e == nil {
		t.Fatal("wrong password accepted")
	}
	if e = s.UpdateScope(context.Background(), id, "云南/昆明/西山/马街"); e != nil {
		t.Fatal(e)
	}
}
func TestIntakePointStore(t *testing.T) {
	db := identityDB(t)
	defer db.Close()
	s := IntakePointStore{DB: db}
	id, e := s.Create(context.Background(), "碧鸡广场", "云南/昆明/西山/马街")
	if e != nil {
		t.Fatal(e)
	}
	x, e := s.Get(context.Background(), id)
	if e != nil || !x.Active {
		t.Fatalf("point=%+v err=%v", x, e)
	}
	if e = s.SetActive(context.Background(), id, false); e != nil {
		t.Fatal(e)
	}
	x, e = s.Get(context.Background(), id)
	if e != nil || x.Active {
		t.Fatal("point remains active")
	}
}
func TestDepartmentStore(t *testing.T) {
	db := identityDB(t)
	defer db.Close()
	s := DepartmentStore{DB: db}
	id, e := s.Create(context.Background(), "区级部门", "云南/昆明/西山/马街")
	if e != nil {
		t.Fatal(e)
	}
	items, e := s.ActiveForScope(context.Background(), "云南/昆明/西山/马街")
	if e != nil || len(items) != 1 {
		t.Fatalf("items=%d err=%v", len(items), e)
	}
	if e = s.Disable(context.Background(), id); e != nil {
		t.Fatal(e)
	}
	items, e = s.ActiveForScope(context.Background(), "云南/昆明/西山/马街")
	if e != nil || len(items) != 0 {
		t.Fatalf("disabled items=%d err=%v", len(items), e)
	}
}
