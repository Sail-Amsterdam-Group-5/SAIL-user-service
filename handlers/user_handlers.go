package handlers

import (
	"encoding/json"
	"net/http"

	"SAIL-user-service/config"
	"SAIL-user-service/keycloak"
	"SAIL-user-service/models"

	"github.com/gorilla/mux"
)

func RegisterUserHandlers(router *mux.Router, cfg *config.Config) {
	client := keycloak.NewKeycloakClient(cfg)

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		users, err := client.GetAllUsers()
		if err != nil {
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(users)
	}).Methods("GET")

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		user, err := client.GetUserById(id)
		if err != nil {
			http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		err := client.RegisterUser(newUser)
		if err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")
}