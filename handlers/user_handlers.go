package handlers

import (
	"encoding/json"
	"net/http"

	"SAIL-user-service/config"
	"SAIL-user-service/keycloak"

	"github.com/gorilla/mux"
)

// @Summary Get all users
// @Description Retrieves all users from the Keycloak server
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Failed to fetch users"
// @Router /users [get]
func GetAllUsersHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := keycloak.NewKeycloakClient(cfg)
		users, err := client.GetAllUsers()
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}
}

// @Summary Get user by ID
// @Description Retrieves a single user from the Keycloak server by userID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 500 {string} string "Failed to fetch user"
// @Router /users/{id} [get]
func GetUserByIDHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		client := keycloak.NewKeycloakClient(cfg)
		id := mux.Vars(r)["id"]
		user, err := client.GetUserById(id)
		if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	}
}

// RegisterUserHandlers registers all user-related handlers
func RegisterUserHandlers(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/users", GetAllUsersHandler(cfg)).Methods("GET")
	router.HandleFunc("/users/{id}", GetUserByIDHandler(cfg)).Methods("GET")
}
