package userstore

import (
	"database/sql"
	"fmt"
	"context"
	_ "github.com/mattn/go-sqlite3"
)

type sqlStore struct {
	db *sql.DB
}

func NewDb(dbPath string) (Store, error) {
	// create sqlite db
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database : %w", err)
	}
	// PRAGMA is sqlite settings
	pragmas := []string{
		// journal_mode - WAL : write ahead logging
		// db changes first write in WAL files and then commit to db
		// its persistence and have better concurrent read/write
		"PRAGMA	journal_mode = WAL;",
		// synchronous settings used in wal mode
		// synchronous controls the fsync operations
		"PRAGMA synchronous = NORMAL;",
		// by default sqlite does not check foreign_keys
		// with this settings it does
		"PRAGMA foreign_keys = ON;",
		// if the db is lock it waits for 5 sec 
		// its good for concurrency and prevent the database is locked error
		//  so it is future-proof
		"PRAGMA busy_timeout=5000;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return nil, fmt.Errorf("failed to apply pragma %s: %w", p, err)
		}
	}

	s := &sqlStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *sqlStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP	
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *sqlStore) Close() error {
	return s.db.Close()
}

// temporary methods to compile the project without error 
// will be implemented later 
func (s *sqlStore) Create(ctx context.Context, user *User) error { return nil }
func (s *sqlStore) GetById(ctx context.Context, id int64) (*User, error) { return nil, nil }
func (s *sqlStore) ListAll(ctx context.Context) ([]User, error) { return nil, nil }
func (s *sqlStore) Update(ctx context.Context, user *User) error { return nil }
func (s *sqlStore) Delete(ctx context.Context, id int64) error { return nil }
