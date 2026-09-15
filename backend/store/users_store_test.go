package store

import (
	"testing"

	"lottery-pool-manager/services"
)

func TestStore_CreateUser(t *testing.T) {
	cleanup(t)

	hash, _ := services.NewAuthService().HashPassword("password123")
	u, err := testStore.CreateUser("Alice", "alice@example.com", hash, "user")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if u.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if u.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", u.Name)
	}
	if u.Email != "alice@example.com" {
		t.Errorf("Email = %v, want alice@example.com", u.Email)
	}
	if u.Role != "user" {
		t.Errorf("Role = %v, want user", u.Role)
	}
	if u.Active {
		t.Error("expected Active to be false by default")
	}
}

func TestStore_CreateUser_DuplicateEmail(t *testing.T) {
	cleanup(t)

	seedUser(t, "Alice", "alice@example.com")

	hash, _ := services.NewAuthService().HashPassword("password123")
	_, err := testStore.CreateUser("Bob", "alice@example.com", hash, "user")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestStore_GetUserByEmail(t *testing.T) {
	cleanup(t)

	created := seedUser(t, "Alice", "alice@example.com")

	u, err := testStore.GetUserByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if u.ID != created.ID {
		t.Errorf("ID = %v, want %v", u.ID, created.ID)
	}
	if u.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", u.Name)
	}
}

func TestStore_GetUserByEmail_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_GetUserByID(t *testing.T) {
	cleanup(t)

	created := seedUser(t, "Alice", "alice@example.com")

	u, err := testStore.GetUserByID(created.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if u.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", u.Name)
	}
}

func TestStore_GetUserByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetUserByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_ActivateUser(t *testing.T) {
	cleanup(t)

	user := seedUser(t, "Alice", "alice@example.com")

	err := testStore.ActivateUser(user.ID)
	if err != nil {
		t.Fatalf("ActivateUser failed: %v", err)
	}

	u, _ := testStore.GetUserByID(user.ID)
	if !u.Active {
		t.Error("expected Active to be true after activation")
	}
}

func TestStore_ActivateUser_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.ActivateUser(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_DeactivateUser(t *testing.T) {
	cleanup(t)

	user := seedUser(t, "Alice", "alice@example.com")

	testStore.ActivateUser(user.ID)
	err := testStore.DeactivateUser(user.ID)
	if err != nil {
		t.Fatalf("DeactivateUser failed: %v", err)
	}

	u, _ := testStore.GetUserByID(user.ID)
	if u.Active {
		t.Error("expected Active to be false after deactivation")
	}
}

func TestStore_DeactivateUser_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.DeactivateUser(999)
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_UpdateUser(t *testing.T) {
	cleanup(t)

	user := seedUser(t, "Alice", "alice@example.com")

	u, err := testStore.UpdateUser(user.ID, "Alice Updated", "alice_new@example.com", "admin")
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if u.Name != "Alice Updated" {
		t.Errorf("Name = %v, want Alice Updated", u.Name)
	}
	if u.Email != "alice_new@example.com" {
		t.Errorf("Email = %v, want alice_new@example.com", u.Email)
	}
	if u.Role != "admin" {
		t.Errorf("Role = %v, want admin", u.Role)
	}
}

func TestStore_UpdateUser_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.UpdateUser(999, "Ghost", "ghost@example.com", "user")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestStore_GetAllUsers(t *testing.T) {
	cleanup(t)

	seedUser(t, "Alice", "alice@example.com")
	seedUser(t, "Bob", "bob@example.com")

	users, err := testStore.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users, got %d", len(users))
	}
}

func TestStore_GetAllUsers_Empty(t *testing.T) {
	cleanup(t)

	users, err := testStore.GetAllUsers()
	if err != nil {
		t.Fatalf("GetAllUsers failed: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("expected 0 users, got %d", len(users))
	}
}
