package userstore

import (
	"context"
	"testing"
)

func StoreTest(t *testing.T) Store {
	t.Helper()

	store, err := NewDb(":memory:")
	if err != nil {
		t.Fatalf("Create DB: %v", err)
	}

	t.Cleanup(func() {
		_ = store.Close()
	})
	return store
}

// Create user test
func TestCreateUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	u := &User{Username: "t", Email: "tst@test.com"}
	err := store.Create(ctx, u)
	if err != nil {
		t.Fatalf("Create failed : %v", err)
	}

	if u.ID == 0 {
		t.Fatal("Expected id to be set got 0")
	}
}

func TestCreateDuplicateUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	u1 := &User{Username: "t", Email: "t@t.com"}
	u2 := &User{Username: "t", Email: "t@t.com"}

	_ = store.Create(ctx, u1)
	err := store.Create(ctx, u2)
	if err != ErrDuplicateUser {
		t.Fatalf("Expected Error from duplicate user but got %v", err)
	}
}

// Get by id test
func TestGetByID(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	u := &User{Username: "t", Email: "t@test.com"}
	_ = store.Create(ctx, u)

	got, err := store.GetById(ctx, u.ID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}
	if got.Username != "t" {
		t.Errorf("Expected t, got %s", got.Username)
	}
}

// user not found test
func TestGetUserNotFound(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	_, err := store.GetById(ctx, 999)
	if err != ErrUserNotFound {
		t.Fatalf("Expected error user not found but got %v", err)
	}
}

// List all test
func TestListAllUsers(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	_ = store.Create(ctx, &User{Username: "t1", Email: "t1@test.com"})
	_ = store.Create(ctx, &User{Username: "t2", Email: "t2@test.com"})

	users, err := store.ListAll(ctx)
	if err != nil {
		t.Fatalf("List failed : %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("Expected 2 users, got %d", len(users))
	}
}

// update test
func TestUpdateUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	u := &User{Username: "a", Email: "a@test.com"}
	_ = store.Create(ctx, u)

	u.Username = "updated"
	err := store.Update(ctx, u)
	if err != nil {
		t.Fatalf("Update failed : %v", err)
	}
}

func TestUpdateNonExistUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	err := store.Update(ctx, &User{
		ID: 999,
		Username: "x",
		Email: "x@test.com",
	})

	if err != ErrUserNotFound{
		t.Fatalf("Expected error user not found got %v", err)
	}
}

// Delete test
func TestDeleteUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	u := &User{
		Username: "t",
		Email: "t@testing.com",
	}
	_ = store.Create(ctx, u)

	err := store.Delete(ctx, u.ID)
	if err != nil {
		t.Fatalf("Delete failed : %v", err)
	}

	_, err = store.GetById(ctx, u.ID)
	if err != ErrUserNotFound {
		t.Fatal("Expected user to be deleted")
	}
}

func TestDeleteNonExistUser(t *testing.T) {
	store := StoreTest(t)
	ctx := context.Background()

	err := store.Delete(ctx, 999)
	if err != ErrUserNotFound {
		t.Fatalf("Expected error user not found, got %v", err)
	}
}

// Close db test
func TestStoreClose(t *testing.T) {
	store := StoreTest(t)
	if err := store.Close(); err != nil {
		t.Errorf("Failed to close store : %v", err)
	}
}

// Empty list test
func TestListEmpty(t *testing.T) {
	store := StoreTest(t)
	users, err := store.ListAll(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users but got %d", len(users))
	}
}

// Detail test
func TestGetByIDDetails(t *testing.T) {
    store := StoreTest(t)
    ctx := context.Background()
    u := &User{Username: "detail_test", Email: "detail@test.com"}
    _ = store.Create(ctx, u)

    got, _ := store.GetById(ctx, u.ID)
    if got.Username != u.Username || got.Email != u.Email {
        t.Error("User details mismatch")
    }
    if got.CreatedAt.IsZero() {
        t.Error("Timestamp should not be zero")
    }
}

// DB test wrong path
func TestNewDbError(t *testing.T){
	_, err := NewDb(".")
	if err == nil {
		t.Error("Expected error for invalid db path, got nil")
	}
}

// Closed db test 
func TestOperationsOnClosedDB(t *testing.T){
	store := StoreTest(t)
	// close db intentionallly
	store.Close()

	ctx := context.Background()
	u := &User{Username: "test", Email: "t@t.com"}

	if err := store.Create(ctx, u); err == nil {
		t.Error("Expected error on closed db for create")
	}

	if _, err := store.ListAll(ctx); err == nil {
		t.Error("Expected error on closed db for ListAll")
	}

	if _, err := store.GetById(ctx, 1); err == nil {
		t.Error("Expected error on closed db for GetById")
	}

	if err := store.Update(ctx, u); err == nil {
		t.Error("Expected error on closed db for Update")
	}

	if err := store.Delete(ctx, 1); err == nil {
		t.Error("Expected error on closed db for Delete")
	}

}
