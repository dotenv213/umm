package userstore

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
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

// CRUD 
func (s *sqlStore) Create(ctx context.Context, user *User) error {
	// Using transactions to make sure it is durable
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Failed to begin transctions : %w", err)
	}
	// if there are some issues in transactions
	// it will rollback the transaction
	defer tx.Rollback()

	// using ? to prevent sql injection from user.
	query := `INSERT INTO users (username, email) VALUES (?, ?)`
	result, err := tx.ExecContext(ctx, query, user.Username, user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed"){
			return ErrDuplicateUser
		}
		return fmt.Errorf("failed to insert user: %w", err)
	}
	// find last id to fill the user struct
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get the last insert id : %w", err)
	}
	user.ID = id

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction : %w", err)
	}
	return nil
}
func (s *sqlStore) GetById(ctx context.Context, id int64) (*User, error) {
	var user User
	query := `SELECT id, username, email, created_at FROM users WHERE id = ?`
	
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}
	return &user, nil
}
func (s *sqlStore) ListAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, username, email, created_at FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list users : %w", err)
	}
	// close rows to free database connection
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user : %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration : %w", err)
	}
	return users, nil
}
func (s *sqlStore) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET username = ?, email = ? WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, user.Username, user.Email, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user : %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotFound
	}
	return nil
}
func (s *sqlStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user : %w", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUserNotFound
	}

	return nil
}
