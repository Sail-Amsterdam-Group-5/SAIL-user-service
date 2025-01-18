package keycloak

import (
	"SAIL-user-service/config"
	"SAIL-user-service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"log"
)

type Client struct {
	config *config.Config
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewKeycloakClient(cfg *config.Config) *Client {
	return &Client{config: cfg}
}

func (kc *Client) GetClientToken() (string, error) {
	log.Println("Getting client token")

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.config.KeycloakURL, kc.config.KeycloakRealm)

	data := "client_id=user-microservice&client_secret=Z7njfSA8YR7kDkftQMKjlqzwM1yqnKLK&grant_type=client_credentials" //GET RID OF HARDCODED SECRET
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		log.Println("Failed to create request")
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to get token")
		return "", fmt.Errorf("Failed to get token: %s", resp.Status)
	}
	log.Println("Got token")

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Println("Failed to decode token")
		return "", err
	}
	log.Println("Succesfully decoded token")

	return tokenResponse.AccessToken, nil
}

func (kc *Client) GetAllUsers() ([]models.User, error) {
	log.Println("Getting all users")

	token, err := kc.GetClientToken()
	if err != nil {
		log.Println("Failed to get token")
		return nil, fmt.Errorf("Failed to get token: %s", err)
	}
	log.Println("Got token")


	url := fmt.Sprintf("%s/admin/realms/%s/users", kc.config.KeycloakURL, kc.config.KeycloakRealm)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create request")
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	log.Println("Created request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Failed to get users")
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Got users")

	var keycloakUsers []models.KeycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUsers); err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(keycloakUsers))
	for _, keycloakUser := range keycloakUsers {
		//get user roles
		roles, err := kc.GetUserRoles(keycloakUser.ID)
		if err != nil {
			return nil, err
		}

		user := keycloakUser.ToUser()
		user.Roles = roles
		users = append(users, *user)
	}
	log.Println("Returning users")

	return users, nil
}

func (kc *Client) GetUserById(id string) (*models.User, error) {
	log.Println("Getting user by id")

	token, err := kc.GetClientToken()
	if err != nil {
		log.Println("Failed to get token")
		return nil, fmt.Errorf("Failed to get token: %s", err)
	}
	log.Println("Got token")

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.config.KeycloakURL, kc.config.KeycloakRealm, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Failed to create request")
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	log.Println("Created request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Failed to get user")
		return nil, err
	}
	defer resp.Body.Close()
	log.Println("Got user")

	var keycloakUser models.KeycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUser); err != nil {
		return nil, err
	}

	//adding the roles
	roles, err := kc.GetUserRoles(id)
	if err != nil {
		log.Println("Failed to get roles")
		return nil, err
	}
	log.Println("Got roles")

	user := keycloakUser.ToUser()
	user.Roles = roles

	log.Println("Returning user")
	
	return user, nil
}

// adding roles to the json bc the /users endpoint doesn't by default
// do this in O(1) lmao
func (kc *Client) GetUserRoles(id string) ([]string, error) {
	token, err := kc.GetClientToken()
	if err != nil {
		return nil, fmt.Errorf("Failed to get token: %s", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/role-mappings/realm", kc.config.KeycloakURL, kc.config.KeycloakRealm, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch roles: %s", resp.Status)
	}

	var keycloakRoles []models.KeycloakRole
	if err := json.NewDecoder(resp.Body).Decode(&keycloakRoles); err != nil {
		return nil, err
	}

	roles := make([]string, len(keycloakRoles))
	for i, role := range keycloakRoles {
		roles[i] = role.Name
	}
	return roles, nil
}
