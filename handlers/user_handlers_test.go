package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// "SAIL-user-service/config"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"SAIL-user-service/keycloak"
	"SAIL-user-service/models"
)

//have to mock the client, so strictly testing handler and not integration with client

func TestGetAllUsers(t *testing.T) {
	mockClient := new(keycloak.MockKeycloakClient) //inshallah
	mockClient.On("GetAllUsers").Return([]models.User{{ID: "1", Username: "Test User"}}, nil)

	// cfg := &config.Config{}
	router := mux.NewRouter()

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) { //have to make it so it uses the mock client
		users, err := mockClient.GetAllUsers()
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var users []models.User
	json.NewDecoder(rr.Body).Decode(&users)
	assert.Equal(t, "Test User", users[0].Username)
	mockClient.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	mockClient := new(keycloak.MockKeycloakClient)
	mockClient.On("GetUserById", "1").Return(&models.User{ID: "1", Username: "Test User"}, nil)

	// cfg := &config.Config{}
	router := mux.NewRouter()

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user, err := mockClient.GetUserById(id)
		if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	}).Methods("GET")

	req, _ := http.NewRequest("GET", "/users/1", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var user models.User
	json.NewDecoder(rr.Body).Decode(&user)
	assert.Equal(t, "Test User", user.Username)
	mockClient.AssertExpectations(t)
}