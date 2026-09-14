package tests

import (
	"time"

	"lottery-pool-manager/models"
)

type mockLoteriaAPIStore struct {
	getLotteryGameByNameFn func(name string) (models.LotteryGame, error)
	createDrawFn           func(gameID int, drawDate time.Time) (models.Draw, error)
	getDrawsByGameFn       func(gameID int) ([]models.Draw, error)
	updateDrawResultsFn    func(id int, resultNumbers *string, resultStars *string) (models.Draw, error)
	updateDrawIDAPIFn      func(id int, drawIDAPI string) (models.Draw, error)
}

func (m *mockLoteriaAPIStore) GetLotteryGameByName(name string) (models.LotteryGame, error) {
	return m.getLotteryGameByNameFn(name)
}
func (m *mockLoteriaAPIStore) CreateDraw(gameID int, drawDate time.Time) (models.Draw, error) {
	return m.createDrawFn(gameID, drawDate)
}
func (m *mockLoteriaAPIStore) GetDrawsByGame(gameID int) ([]models.Draw, error) {
	return m.getDrawsByGameFn(gameID)
}
func (m *mockLoteriaAPIStore) UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
	return m.updateDrawResultsFn(id, resultNumbers, resultStars)
}
func (m *mockLoteriaAPIStore) UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error) {
	return m.updateDrawIDAPIFn(id, drawIDAPI)
}
