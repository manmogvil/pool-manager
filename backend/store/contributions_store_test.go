package store

import (
	"testing"

	"lottery-pool-manager/models"
)

func TestStore_CreateContribution(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	c := &models.Contribution{
		ParticipantID: participant.ID,
		GameID:        game.ID,
		Month:         1,
		Year:          2024,
		Amount:        5.00,
		Paid:          true,
		PaymentMethod: "CASH",
	}

	result, err := testStore.CreateContribution(c)
	if err != nil {
		t.Fatalf("CreateContribution failed: %v", err)
	}
	if result.ID == 0 {
		t.Error("expected non-zero ID")
	}
	if result.Amount != 5.00 {
		t.Errorf("Amount = %v, want 5.00", result.Amount)
	}
}

func TestStore_GetAllContributions(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})
	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 2, Year: 2024, Amount: 5.00, PaymentMethod: "BIZUM",
	})

	contributions, err := testStore.GetAllContributions()
	if err != nil {
		t.Fatalf("GetAllContributions failed: %v", err)
	}
	if len(contributions) != 2 {
		t.Errorf("expected 2 contributions, got %d", len(contributions))
	}
}

func TestStore_GetContributionByID(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	created, _ := testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})

	c, err := testStore.GetContributionByID(created.ID)
	if err != nil {
		t.Fatalf("GetContributionByID failed: %v", err)
	}
	if c.Amount != 5.00 {
		t.Errorf("Amount = %v, want 5.00", c.Amount)
	}
}

func TestStore_GetContributionByID_NotFound(t *testing.T) {
	cleanup(t)

	_, err := testStore.GetContributionByID(999)
	if err == nil {
		t.Fatal("expected error for non-existent contribution")
	}
}

func TestStore_GetContributionsByParticipant(t *testing.T) {
	cleanup(t)

	p1 := seedParticipant(t, "Alice", "alice@example.com")
	p2 := seedParticipant(t, "Bob", "bob@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: p1.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})
	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: p1.ID, GameID: game.ID, Month: 2, Year: 2024, Amount: 5.00, PaymentMethod: "BIZUM",
	})
	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: p2.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})

	contributions, err := testStore.GetContributionsByParticipant(p1.ID)
	if err != nil {
		t.Fatalf("GetContributionsByParticipant failed: %v", err)
	}
	if len(contributions) != 2 {
		t.Errorf("expected 2 contributions for p1, got %d", len(contributions))
	}
}

func TestStore_GetContributionsByPeriod(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})
	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 2, Year: 2024, Amount: 5.00, PaymentMethod: "BIZUM",
	})
	_, _ = testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2023, Amount: 5.00, PaymentMethod: "CASH",
	})

	contributions, err := testStore.GetContributionsByPeriod(1, 2024)
	if err != nil {
		t.Fatalf("GetContributionsByPeriod failed: %v", err)
	}
	if len(contributions) != 1 {
		t.Errorf("expected 1 contribution for period 1/2024, got %d", len(contributions))
	}
}

func TestStore_UpdateContribution(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	created, _ := testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})

	updated := &models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 10.00, Paid: true, PaymentMethod: "BIZUM",
	}

	result, err := testStore.UpdateContribution(created.ID, updated)
	if err != nil {
		t.Fatalf("UpdateContribution failed: %v", err)
	}
	if result.Amount != 10.00 {
		t.Errorf("Amount = %v, want 10.00", result.Amount)
	}
}

func TestStore_UpdateContribution_NotFound(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	updated := &models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 10.00, PaymentMethod: "BIZUM",
	}

	_, err := testStore.UpdateContribution(999, updated)
	if err == nil {
		t.Fatal("expected error for non-existent contribution")
	}
}

func TestStore_DeleteContribution(t *testing.T) {
	cleanup(t)

	participant := seedParticipant(t, "Alice", "alice@example.com")
	game := seedGame(t, "EuroMillones", "Tuesday, Friday", 2.50)

	created, _ := testStore.CreateContribution(&models.Contribution{
		ParticipantID: participant.ID, GameID: game.ID, Month: 1, Year: 2024, Amount: 5.00, PaymentMethod: "CASH",
	})

	err := testStore.DeleteContribution(created.ID)
	if err != nil {
		t.Fatalf("DeleteContribution failed: %v", err)
	}

	_, err = testStore.GetContributionByID(created.ID)
	if err == nil {
		t.Fatal("expected error for deleted contribution")
	}
}

func TestStore_DeleteContribution_NotFound(t *testing.T) {
	cleanup(t)

	err := testStore.DeleteContribution(999)
	if err == nil {
		t.Fatal("expected error for non-existent contribution")
	}
}
