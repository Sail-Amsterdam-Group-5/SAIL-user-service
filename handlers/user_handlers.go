package handlers

import (
	"encoding/json"
	"net/http"

	"SAIL-user-service/config"
	"SAIL-user-service/keycloak"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// @Summary Get all users
// @Description Retrieves all users from the Keycloak server
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Failed to fetch users"
// @Router /users [get]
func GetAllUsersHandler(cfg *config.Config) http.HandlerFunc {
	logger := logrus.WithFields(logrus.Fields{
		"handler": "GetAllUsersHandler",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request to fetch all users")

		client := keycloak.NewKeycloakClient(cfg)
		users, err := client.GetAllUsers()
		if err != nil {
			logger.WithError(err).Error("Failed to fetch users")
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		
		logger.WithField("user_count", len(users)).Info("Successfully fetched all users")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			logger.WithError(err).Error("Failed to encode users response")
			http.Error(w, "Failed to process response", http.StatusInternalServerError)
		}
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
	logger := logrus.WithFields(logrus.Fields{
		"handler": "GetUserByIDHandler",
	})

	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		logger.WithField("user_id", id).Info("Received request to fetch user by ID")

		client := keycloak.NewKeycloakClient(cfg)

		user, err := client.GetUserById(id)
		if err != nil {
			logger.WithError(err).WithField("user_id", id).Error("Failed to fetch user")
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
		
		logger.WithField("user_id", id).Info("Successfully fetched user")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logger.WithError(err).WithField("user_id", id).Error("Failed to encode user response")
			http.Error(w, "Failed to process response", http.StatusInternalServerError)
		}
	}
}

// RegisterUserHandlers registers all user-related handlers
func RegisterUserHandlers(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/users", GetAllUsersHandler(cfg)).Methods("GET")
	router.HandleFunc("/users/{id}", GetUserByIDHandler(cfg)).Methods("GET")
}
