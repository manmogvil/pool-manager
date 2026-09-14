package tests

import (
	"lottery-pool-manager/models"
)

type mockCheckTicketStore struct {
	getTicketByIDFn      func(id int) (models.Ticket, error)
	getDrawByIDFn        func(id int) (models.Draw, error)
	getLotteryGameByIDFn func(id int) (models.LotteryGame, error)
	updateTicketPrizeFn  func(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error)
}

func (m *mockCheckTicketStore) GetTicketByID(id int) (models.Ticket, error) {
	return m.getTicketByIDFn(id)
}
func (m *mockCheckTicketStore) GetDrawByID(id int) (models.Draw, error) {
	return m.getDrawByIDFn(id)
}
func (m *mockCheckTicketStore) GetLotteryGameByID(id int) (models.LotteryGame, error) {
	return m.getLotteryGameByIDFn(id)
}
func (m *mockCheckTicketStore) UpdateTicketPrize(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error) {
	return m.updateTicketPrizeFn(id, prizeTier, prizeAmount, matchedNumbers, matchedStars)
}
