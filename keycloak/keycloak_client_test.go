package keycloak

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"SAIL-user-service/config"
	"SAIL-user-service/models"

	"github.com/stretchr/testify/assert"
)

func TestGetClientToken(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { //just PLEASE
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/realms/mock-realm/protocol/openid-connect/token", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"access_token": "mock-token"}`))
	}))
	defer mockServer.Close()

	client := NewKeycloakClient(&config.Config{
		KeycloakURL:  mockServer.URL,
		KeycloakRealm: "mock-realm",
	})

	t.Run("Success", func(t *testing.T) {
		token, err := client.GetClientToken()
		assert.NoError(t, err)
		assert.Equal(t, "mock-token", token)
	})

	t.Run("Failed request", func(t *testing.T) {
		client := NewKeycloakClient(&config.Config{
			KeycloakURL:  "http://invalid-url",
			KeycloakRealm: "mock-realm",
		})

		token, err := client.GetClientToken()
		assert.Error(t, err)
		assert.Empty(t, token)
	})

	t.Run("Failed server", func(t *testing.T) {
		mockServer.Close()

		token, err := client.GetClientToken()
		assert.Error(t, err)
		assert.Empty(t, token)
	})
}

func TestGetAllUsers(t *testing.T) {
	mockClient := new(MockKeycloakClient)

	mockUsers := []models.User{
		{ID: "1", Username: "user1"},
		{ID: "2", Username: "user2"},
	}

	mockClient.On("GetAllUsers").Return(mockUsers, nil)

	users, err := mockClient.GetAllUsers()
	assert.NoError(t, err)
	assert.Equal(t, mockUsers, users)

	mockClient.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	mockClient := new(MockKeycloakClient)
	mockUser := &models.User{ID: "1", Username: "user1"}

	mockClient.On("GetUserById", "1").Return(mockUser, nil)

	user, err := mockClient.GetUserById("1")
	assert.NoError(t, err)
	assert.Equal(t, mockUser, user)

	mockClient.AssertExpectations(t)
}

func TestGetUserRoles(t *testing.T) {
	mockClient := new(MockKeycloakClient)
	mockRoles := []string{"admin", "user"}

	mockClient.On("GetUserRoles", "1").Return(mockRoles, nil)

	roles, err := mockClient.GetUserRoles("1")
	assert.NoError(t, err)
	assert.Equal(t, mockRoles, roles)

	mockClient.AssertExpectations(t)
}
