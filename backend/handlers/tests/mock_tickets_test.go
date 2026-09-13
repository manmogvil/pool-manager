package tests

import (
	"lottery-pool-manager/models"
)

type mockTicketStore struct {
	createFn      func(drawID int, numbers string, stars *string, cost float64) (models.Ticket, error)
	getAllFn      func() ([]models.Ticket, error)
	getByIDFn     func(id int) (models.Ticket, error)
	getByDrawFn   func(drawID int) ([]models.Ticket, error)
	updatePrizeFn func(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error)
	deleteFn      func(id int) error
}

func (m *mockTicketStore) CreateTicket(drawID int, numbers string, stars *string, cost float64) (models.Ticket, error) {
	return m.createFn(drawID, numbers, stars, cost)
}
func (m *mockTicketStore) GetAllTickets() ([]models.Ticket, error) {
	return m.getAllFn()
}
func (m *mockTicketStore) GetTicketByID(id int) (models.Ticket, error) {
	return m.getByIDFn(id)
}
func (m *mockTicketStore) GetTicketsByDraw(drawID int) ([]models.Ticket, error) {
	return m.getByDrawFn(drawID)
}
func (m *mockTicketStore) UpdateTicketPrize(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error) {
	return m.updatePrizeFn(id, prizeTier, prizeAmount, matchedNumbers, matchedStars)
}
func (m *mockTicketStore) DeleteTicket(id int) error {
	return m.deleteFn(id)
}
