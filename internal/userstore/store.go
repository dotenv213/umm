package userstore

import "context"

// This interface is a contract that
// represent how crud implemented in this module
type Store interface {
	Create(ctx context.Context, user *User) error
	GetById(ctx context.Context, id int64) (*User, error)
	ListAll(ctx context.Context)([]User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Close() error	
}