package store

import (
	"testing"
	"time"
)

func TestStore_CreateTicket(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	ticket, err := testStore.CreateTicket(draw.ID, "5,12,23,34,45", "3,7", 2.50)
	if err != nil {
		t.Fatalf("CreateTicket failed: %v", err)
	}
	if ticket.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if ticket.Numbers != "5,12,23,34,45" {
		t.Errorf("Numbers = %v, want 5,12,23,34,45", ticket.Numbers)
	}
}

func TestStore_GetAllTickets(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	_, _ = testStore.CreateTicket(draw.ID, "1,2,3,4,5", "1,2", 2.50)
	_, _ = testStore.CreateTicket(draw.ID, "6,7,8,9,10", "3,4", 2.50)

	tickets, err := testStore.GetAllTickets()
	if err != nil {
		t.Fatalf("GetAllTickets failed: %v", err)
	}
	if len(tickets) != 2 {
		t.Errorf("expected 2 tickets, got %d", len(tickets))
	}
}

func TestStore_GetTicketByID(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	created, _ := testStore.CreateTicket(draw.ID, "5,12,23,34,45", "3,7", 2.50)

	ticket, err := testStore.GetTicketByID(created.ID)
	if err != nil {
		t.Fatalf("GetTicketByID failed: %v", err)
	}
	if ticket.Numbers != "5,12,23,34,45" {
		t.Errorf("Numbers = %v, want 5,12,23,34,45", ticket.Numbers)
	}
}

func TestStore_GetTicketByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetTicketByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent ticket")
	}
}

func TestStore_GetTicketsByDraw(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw1 := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))
	draw2 := seedDraw(t, game.ID, time.Date(2024, 1, 22, 0, 0, 0, 0, time.UTC))

	_, _ = testStore.CreateTicket(draw1.ID, "1,2,3,4,5", "1,2", 2.50)
	_, _ = testStore.CreateTicket(draw1.ID, "6,7,8,9,10", "3,4", 2.50)
	_, _ = testStore.CreateTicket(draw2.ID, "11,12,13,14,15", "5,6", 2.50)

	tickets, err := testStore.GetTicketsByDraw(draw1.ID)
	if err != nil {
		t.Fatalf("GetTicketsByDraw failed: %v", err)
	}
	if len(tickets) != 2 {
		t.Errorf("expected 2 tickets for draw1, got %d", len(tickets))
	}
}

func TestStore_UpdateTicketPrize(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	created, _ := testStore.CreateTicket(draw.ID, "5,12,23,34,45", "3,7", 2.50)

	prizeTier := "5th"
	prizeAmount := 50.0
	matchedNumbers := 3
	matchedStars := 1

	ticket, err := testStore.UpdateTicketPrize(created.ID, &prizeTier, &prizeAmount, &matchedNumbers, &matchedStars)
	if err != nil {
		t.Fatalf("UpdateTicketPrize failed: %v", err)
	}
	if *ticket.PrizeTier != "5th" {
		t.Errorf("PrizeTier = %v, want 5th", *ticket.PrizeTier)
	}
	if *ticket.PrizeAmount != 50.0 {
		t.Errorf("PrizeAmount = %v, want 50.0", *ticket.PrizeAmount)
	}
}

func TestStore_UpdateTicketPrize_NotFound(t *testing.T) {
	cleanup(t)

	prizeTier := "5th"
	_, err := testStore.UpdateTicketPrize(999, &prizeTier, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for non-existent ticket")
	}
}

func TestStore_DeleteTicket(t *testing.T) {
	cleanup(t)

	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)
	draw := seedDraw(t, game.ID, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC))

	created, _ := testStore.CreateTicket(draw.ID, "5,12,23,34,45", "3,7", 2.50)

	err := testStore.DeleteTicket(created.ID)
	if err != nil {
		t.Fatalf("DeleteTicket failed: %v", err)
	}

	_, err = testStore.GetTicketByID(created.ID)
	if err == nil {
		t.Fatal("expected error for deleted ticket")
	}
}

func TestStore_DeleteTicket_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.DeleteTicket(999)
	if err == nil {
		t.Fatal("expected error for non-existent ticket")
	}
}
