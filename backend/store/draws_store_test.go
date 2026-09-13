package store

import (
	"testing"
	"time"
)

func TestStore_CreateDraw(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	drawDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	d, err := testStore.CreateDraw(game.ID, drawDate)
	if err != nil {
		t.Fatalf("CreateDraw failed: %v", err)
	}
	if d.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if d.GameID != game.ID {
		t.Errorf("GameID = %v, want %v", d.GameID, game.ID)
	}
}

func TestStore_CreateDraw_Duplicate(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	drawDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	seedDraw(t, game.ID, drawDate)

	_, err := testStore.CreateDraw(game.ID, drawDate)
	if err == nil {
		t.Fatal("expected error for duplicate draw")
	}
}

func TestStore_GetAllDraws(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	seedDraw(t, game.ID, time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC))

	draws, err := testStore.GetAllDraws()
	if err != nil {
		t.Fatalf("GetAllDraws failed: %v", err)
	}
	if len(draws) != 2 {
		t.Errorf("expected 2 draws, got %d", len(draws))
	}
}

func TestStore_GetDrawByID(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	created := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	d, err := testStore.GetDrawByID(created.ID)
	if err != nil {
		t.Fatalf("GetDrawByID failed: %v", err)
	}
	if d.GameID != game.ID {
		t.Errorf("GameID = %v, want %v", d.GameID, game.ID)
	}
}

func TestStore_GetDrawByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetDrawByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent draw")
	}
}

func TestStore_GetDrawsByGame(t *testing.T) {
	cleanup(t)

	game1 := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	game2 := seedGame(t, "El Gordo", "Thursday", 1.50)

	seedDraw(t, game1.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	seedDraw(t, game1.ID, time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC))
	seedDraw(t, game2.ID, time.Date(2024, 1, 18, 0, 0, 0, 0, time.UTC))

	draws, err := testStore.GetDrawsByGame(game1.ID)
	if err != nil {
		t.Fatalf("GetDrawsByGame failed: %v", err)
	}
	if len(draws) != 2 {
		t.Errorf("expected 2 draws for game1, got %d", len(draws))
	}
}

func TestStore_GetPendingDraws(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	d1 := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	seedDraw(t, game.ID, time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC))

	_ = testStore.MarkDrawAsProcessed(d1.ID)

	draws, err := testStore.GetPendingDraws()
	if err != nil {
		t.Fatalf("GetPendingDraws failed: %v", err)
	}
	if len(draws) != 1 {
		t.Errorf("expected 1 pending draw, got %d", len(draws))
	}
}

func TestStore_UpdateDrawResults(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	created := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	numbers := "5,12,23,34,45"
	stars := "3,7"

	d, err := testStore.UpdateDrawResults(created.ID, &numbers, &stars)
	if err != nil {
		t.Fatalf("UpdateDrawResults failed: %v", err)
	}
	if *d.ResultNumbers != numbers {
		t.Errorf("ResultNumbers = %v, want %v", *d.ResultNumbers, numbers)
	}
	if *d.ResultStars != stars {
		t.Errorf("ResultStars = %v, want %v", *d.ResultStars, stars)
	}
}

func TestStore_UpdateDrawResults_NotFound(t *testing.T) {
	cleanup(t)

	numbers := "5,12,23,34,45"
	_, err := testStore.UpdateDrawResults(999, &numbers, nil)
	if err == nil {
		t.Fatal("expected error for non-existent draw")
	}
}

func TestStore_MarkDrawAsProcessed(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	created := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	err := testStore.MarkDrawAsProcessed(created.ID)
	if err != nil {
		t.Fatalf("MarkDrawAsProcessed failed: %v", err)
	}

	d, _ := testStore.GetDrawByID(created.ID)
	if !d.Processed {
		t.Error("expected Processed to be true")
	}
}

func TestStore_MarkDrawAsProcessed_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.MarkDrawAsProcessed(999)
	if err == nil {
		t.Fatal("expected error for non-existent draw")
	}
}

func TestStore_DeleteDraw(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	created := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	err := testStore.DeleteDraw(created.ID)
	if err != nil {
		t.Fatalf("DeleteDraw failed: %v", err)
	}

	_, err = testStore.GetDrawByID(created.ID)
	if err == nil {
		t.Fatal("expected error for deleted draw")
	}
}

func TestStore_DeleteDraw_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.DeleteDraw(999)
	if err == nil {
		t.Fatal("expected error for non-existent draw")
	}
}

func TestStore_DeleteDraw_HasTickets(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	_, err := testStore.CreateTicket(draw.ID, "1,2,3,4,5", "1,2", 2.50)
	if err != nil {
		t.Fatalf("CreateTicket failed: %v", err)
	}

	err = testStore.DeleteDraw(draw.ID)
	if err == nil {
		t.Fatal("expected error when deleting draw with tickets")
	}
}
