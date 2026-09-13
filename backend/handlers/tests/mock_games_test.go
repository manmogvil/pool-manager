package tests

import (
	"lottery-pool-manager/models"
)

type mockLotteryGameStore struct {
	createFn  func(name, drawDays string, ticketPrice float64) (models.LotteryGame, error)
	getAllFn  func() ([]models.LotteryGame, error)
	getByIDFn func(id int) (models.LotteryGame, error)
	updateFn  func(id int, name, drawDays string, ticketPrice float64, active bool) (models.LotteryGame, error)
	deleteFn  func(id int) error
}

func (m *mockLotteryGameStore) CreateLotteryGame(name, drawDays string, ticketPrice float64) (models.LotteryGame, error) {
	return m.createFn(name, drawDays, ticketPrice)
}
func (m *mockLotteryGameStore) GetAllLotteryGames() ([]models.LotteryGame, error) {
	return m.getAllFn()
}
func (m *mockLotteryGameStore) GetLotteryGameByID(id int) (models.LotteryGame, error) {
	return m.getByIDFn(id)
}
func (m *mockLotteryGameStore) UpdateLotteryGame(id int, name, drawDays string, ticketPrice float64, active bool) (models.LotteryGame, error) {
	return m.updateFn(id, name, drawDays, ticketPrice, active)
}
func (m *mockLotteryGameStore) DeleteLotteryGame(id int) error {
	return m.deleteFn(id)
}
