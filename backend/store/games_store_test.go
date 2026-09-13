package store

import (
	"testing"
	"time"
)

func TestStore_CreateLotteryGame(t *testing.T) {
	cleanup(t)

	g, err := testStore.CreateLotteryGame("EuroMillones", "Tuesday, Friday", 2.50)
	if err != nil {
		t.Fatalf("CreateLotteryGame failed: %v", err)
	}
	if g.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if g.Name != "EuroMillones" {
		t.Errorf("Name = %v, want EuroMillones", g.Name)
	}
	if !g.Active {
		t.Error("expected Active to be true")
	}
}

func TestStore_CreateLotteryGame_DuplicateName(t *testing.T) {
	cleanup(t)

	seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	_, err := testStore.CreateLotteryGame("EuroMillones", "Monday", 3.00)
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestStore_GetAllLotteryGames(t *testing.T) {
	cleanup(t)

	seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	seedGame(t, "El Gordo", "Thursday", 1.50)

	games, err := testStore.GetAllLotteryGames()
	if err != nil {
		t.Fatalf("GetAllLotteryGames failed: %v", err)
	}
	if len(games) != 2 {
		t.Errorf("expected 2 games, got %d", len(games))
	}
}

func TestStore_GetLotteryGameByID(t *testing.T) {
	cleanup(t)

	created := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	g, err := testStore.GetLotteryGameByID(created.ID)
	if err != nil {
		t.Fatalf("GetLotteryGameByID failed: %v", err)
	}
	if g.Name != "EuroMillones" {
		t.Errorf("Name = %v, want EuroMillones", g.Name)
	}
}

func TestStore_GetLotteryGameByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetLotteryGameByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}
}

func TestStore_UpdateLotteryGame(t *testing.T) {
	cleanup(t)

	created := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	g, err := testStore.UpdateLotteryGame(created.ID, "EuroMillones Updated", "Monday, Wednesday", 3.00, false)
	if err != nil {
		t.Fatalf("UpdateLotteryGame failed: %v", err)
	}
	if g.Name != "EuroMillones Updated" {
		t.Errorf("Name = %v, want EuroMillones Updated", g.Name)
	}
	if g.Active {
		t.Error("expected Active to be false")
	}
}

func TestStore_UpdateLotteryGame_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.UpdateLotteryGame(999, "Ghost", "Monday", 1.00, true)
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}
}

func TestStore_DeleteLotteryGame(t *testing.T) {
	cleanup(t)

	created := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	err := testStore.DeleteLotteryGame(created.ID)
	if err != nil {
		t.Fatalf("DeleteLotteryGame failed: %v", err)
	}

	_, err = testStore.GetLotteryGameByID(created.ID)
	if err == nil {
		t.Fatal("expected error for deleted game")
	}
}

func TestStore_DeleteLotteryGame_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.DeleteLotteryGame(999)
	if err == nil {
		t.Fatal("expected error for non-existent game")
	}
}

func TestStore_DeleteLotteryGame_HasTickets(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Now().AddDate(0, 0, 7))

	_, err := testStore.CreateTicket(draw.ID, "1,2,3,4,5", "1,2", 2.50)
	if err != nil {
		t.Fatalf("CreateTicket failed: %v", err)
	}

	err = testStore.DeleteLotteryGame(game.ID)
	if err == nil {
		t.Fatal("expected error when deleting game with tickets")
	}
}
