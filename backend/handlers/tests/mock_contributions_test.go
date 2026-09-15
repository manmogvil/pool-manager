package tests

import (
	"lottery-pool-manager/models"
)

type mockContributionStore struct {
	createFn           func(c *models.Contribution) (models.Contribution, error)
	getAllFn           func() ([]models.Contribution, error)
	getByIDFn          func(id int) (models.Contribution, error)
	getByParticipantFn func(userID int) ([]models.Contribution, error)
	getByPeriodFn      func(month, year int) ([]models.Contribution, error)
	updateFn           func(id int, c *models.Contribution) (models.Contribution, error)
	deleteFn           func(id int) error
}

func (m *mockContributionStore) CreateContribution(c *models.Contribution) (models.Contribution, error) {
	return m.createFn(c)
}
func (m *mockContributionStore) GetAllContributions() ([]models.Contribution, error) {
	return m.getAllFn()
}
func (m *mockContributionStore) GetContributionByID(id int) (models.Contribution, error) {
	return m.getByIDFn(id)
}
func (m *mockContributionStore) GetContributionsByParticipant(userID int) ([]models.Contribution, error) {
	return m.getByParticipantFn(userID)
}
func (m *mockContributionStore) GetContributionsByPeriod(month, year int) ([]models.Contribution, error) {
	return m.getByPeriodFn(month, year)
}
func (m *mockContributionStore) UpdateContribution(id int, c *models.Contribution) (models.Contribution, error) {
	return m.updateFn(id, c)
}
func (m *mockContributionStore) DeleteContribution(id int) error {
	return m.deleteFn(id)
}
