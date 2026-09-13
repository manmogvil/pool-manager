package tests

import (
	"lottery-pool-manager/models"
)

type mockParticipantStore struct {
	getAllFn     func() ([]models.Participant, error)
	getByIDFn    func(id int) (models.Participant, error)
	createFn     func(name, email string) (models.Participant, error)
	updateFn     func(id int, name, email string) (models.Participant, error)
	deactivateFn func(id int) (models.Participant, error)
	activateFn   func(id int) (models.Participant, error)
}

func (m *mockParticipantStore) GetAllParticipants() ([]models.Participant, error) {
	return m.getAllFn()
}
func (m *mockParticipantStore) GetParticipantByID(id int) (models.Participant, error) {
	return m.getByIDFn(id)
}
func (m *mockParticipantStore) CreateParticipant(name, email string) (models.Participant, error) {
	return m.createFn(name, email)
}
func (m *mockParticipantStore) UpdateParticipant(id int, name, email string) (models.Participant, error) {
	return m.updateFn(id, name, email)
}
func (m *mockParticipantStore) DeactivateParticipant(id int) (models.Participant, error) {
	return m.deactivateFn(id)
}
func (m *mockParticipantStore) ActivateParticipant(id int) (models.Participant, error) {
	return m.activateFn(id)
}
