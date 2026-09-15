package tests

import (
	"lottery-pool-manager/models"
)

type mockAuthStore struct {
	createUserFn    func(name, email, passwordHash, role string) (models.User, error)
	getByEmailFn    func(email string) (models.User, error)
	getByIDFn       func(id int) (models.User, error)
	activateFn      func(id int) error
	deactivateFn    func(id int) error
	updateUserFn    func(id int, name, email, role string) (models.User, error)
	getAllUsersFn   func() ([]models.User, error)
}

func (m *mockAuthStore) CreateUser(name, email, passwordHash, role string) (models.User, error) {
	return m.createUserFn(name, email, passwordHash, role)
}
func (m *mockAuthStore) GetUserByEmail(email string) (models.User, error) {
	return m.getByEmailFn(email)
}
func (m *mockAuthStore) GetUserByID(id int) (models.User, error) {
	return m.getByIDFn(id)
}
func (m *mockAuthStore) ActivateUser(id int) error {
	return m.activateFn(id)
}
func (m *mockAuthStore) DeactivateUser(id int) error {
	return m.deactivateFn(id)
}
func (m *mockAuthStore) UpdateUser(id int, name, email, role string) (models.User, error) {
	return m.updateUserFn(id, name, email, role)
}
func (m *mockAuthStore) GetAllUsers() ([]models.User, error) {
	return m.getAllUsersFn()
}
