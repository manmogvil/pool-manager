package store

import (
	"testing"
)

func TestStore_CreateParticipant(t *testing.T) {
	cleanup(t)

	p, err := testStore.CreateParticipant("Alice", "alice@example.com")
	if err != nil {
		t.Fatalf("CreateParticipant failed: %v", err)
	}
	if p.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if p.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", p.Name)
	}
	if p.Email != "alice@example.com" {
		t.Errorf("Email = %v, want alice@example.com", p.Email)
	}
	if !p.Active {
		t.Error("expected Active to be true")
	}
}

func TestStore_CreateParticipant_DuplicateEmail(t *testing.T) {
	cleanup(t)

	seedParticipant(t, "Alice", "alice@example.com")

	_, err := testStore.CreateParticipant("Bob", "alice@example.com")
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestStore_GetAllParticipants(t *testing.T) {
	cleanup(t)

	seedParticipant(t, "Alice", "alice@example.com")
	seedParticipant(t, "Bob", "bob@example.com")

	participants, err := testStore.GetAllParticipants()
	if err != nil {
		t.Fatalf("GetAllParticipants failed: %v", err)
	}
	if len(participants) != 2 {
		t.Errorf("expected 2 participants, got %d", len(participants))
	}
}

func TestStore_GetParticipantByID(t *testing.T) {
	cleanup(t)

	created := seedParticipant(t, "Alice", "alice@example.com")

	p, err := testStore.GetParticipantByID(created.ID)
	if err != nil {
		t.Fatalf("GetParticipantByID failed: %v", err)
	}
	if p.Name != "Alice" {
		t.Errorf("Name = %v, want Alice", p.Name)
	}
}

func TestStore_GetParticipantByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetParticipantByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent participant")
	}
}

func TestStore_UpdateParticipant(t *testing.T) {
	cleanup(t)

	created := seedParticipant(t, "Alice", "alice@example.com")

	p, err := testStore.UpdateParticipant(created.ID, "Alice Updated", "alice_new@example.com")
	if err != nil {
		t.Fatalf("UpdateParticipant failed: %v", err)
	}
	if p.Name != "Alice Updated" {
		t.Errorf("Name = %v, want Alice Updated", p.Name)
	}
	if p.Email != "alice_new@example.com" {
		t.Errorf("Email = %v, want alice_new@example.com", p.Email)
	}
}

func TestStore_UpdateParticipant_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.UpdateParticipant(999, "Ghost", "ghost@example.com")
	if err == nil {
		t.Fatal("expected error for non-existent participant")
	}
}

func TestStore_DeactivateParticipant(t *testing.T) {
	cleanup(t)

	created := seedParticipant(t, "Alice", "alice@example.com")

	p, err := testStore.DeactivateParticipant(created.ID)
	if err != nil {
		t.Fatalf("DeactivateParticipant failed: %v", err)
	}
	if p.Active {
		t.Error("expected Active to be false")
	}
}

func TestStore_ActivateParticipant(t *testing.T) {
	cleanup(t)

	created := seedParticipant(t, "Alice", "alice@example.com")
	_, _ = testStore.DeactivateParticipant(created.ID)

	p, err := testStore.ActivateParticipant(created.ID)
	if err != nil {
		t.Fatalf("ActivateParticipant failed: %v", err)
	}
	if !p.Active {
		t.Error("expected Active to be true")
	}
}
