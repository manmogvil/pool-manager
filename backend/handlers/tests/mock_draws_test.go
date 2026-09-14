package tests

import (
	"time"

	"lottery-pool-manager/models"
)

type mockDrawStore struct {
	createFn            func(gameID int, drawDate time.Time) (models.Draw, error)
	getAllFn            func() ([]models.Draw, error)
	getByIDFn           func(id int) (models.Draw, error)
	getByGameFn         func(gameID int) ([]models.Draw, error)
	getPendingFn        func() ([]models.Draw, error)
	updateResultsFn     func(id int, resultNumbers *string, resultStars *string) (models.Draw, error)
	updateDrawIDAPIFn   func(id int, drawIDAPI string) (models.Draw, error)
	markProcessedFn     func(id int) error
	deleteFn            func(id int) error
}

func (m *mockDrawStore) CreateDraw(gameID int, drawDate time.Time) (models.Draw, error) {
	return m.createFn(gameID, drawDate)
}
func (m *mockDrawStore) GetAllDraws() ([]models.Draw, error) {
	return m.getAllFn()
}
func (m *mockDrawStore) GetDrawByID(id int) (models.Draw, error) {
	return m.getByIDFn(id)
}
func (m *mockDrawStore) GetDrawsByGame(gameID int) ([]models.Draw, error) {
	return m.getByGameFn(gameID)
}
func (m *mockDrawStore) GetPendingDraws() ([]models.Draw, error) {
	return m.getPendingFn()
}
func (m *mockDrawStore) UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
	return m.updateResultsFn(id, resultNumbers, resultStars)
}
func (m *mockDrawStore) UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error) {
	return m.updateDrawIDAPIFn(id, drawIDAPI)
}
func (m *mockDrawStore) MarkDrawAsProcessed(id int) error {
	return m.markProcessedFn(id)
}
func (m *mockDrawStore) DeleteDraw(id int) error {
	return m.deleteFn(id)
}
