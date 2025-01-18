package keycloak

import (
	"SAIL-user-service/models"
	"github.com/stretchr/testify/mock"
)

type MockKeycloakClient struct {
	mock.Mock
}

func (m *MockKeycloakClient) GetClientToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockKeycloakClient) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockKeycloakClient) GetUserById(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockKeycloakClient) GetUserRoles(id string) ([]string, error) {
	args := m.Called(id)
	return args.Get(0).([]string), args.Error(1)
}