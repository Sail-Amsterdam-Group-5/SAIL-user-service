package keycloak

import (
	"SAIL-user-service/config"
	"SAIL-user-service/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Config struct {
	KeycloakURL   string
	KeycloakRealm string
	AdminToken    string
}

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
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.config.KeycloakURL, kc.config.KeycloakRealm)

	data := "client_id=user-microservice&client_secret=Z7njfSA8YR7kDkftQMKjlqzwM1yqnKLK&grant_type=client_credentials"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to get token: %s", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (kc *Client) GetAllUsers() ([]models.User, error) {
	token, err := kc.GetClientToken()
	if err != nil {
		return nil, fmt.Errorf("Failed to get token: %s", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", kc.config.KeycloakURL, kc.config.KeycloakRealm)
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
	return users, nil
}

func (kc *Client) GetUserById(id string) (*models.User, error) {
	token, err := kc.GetClientToken()
	if err != nil {
		return nil, fmt.Errorf("Failed to get token: %s", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.config.KeycloakURL, kc.config.KeycloakRealm, id)
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

	var keycloakUser models.KeycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUser); err != nil {
		return nil, err
	}

	//adding the roles
	roles, err := kc.GetUserRoles(id)
	if err != nil {
		return nil, err
	}

	user := keycloakUser.ToUser()
	user.Roles = roles

	return user, nil
}

func (kc *Client) RegisterUser(newUser models.User) error {
	token, err := kc.GetClientToken()
	if err != nil {
		return fmt.Errorf("Failed to get token: %s", err)
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users", kc.config.KeycloakURL, kc.config.KeycloakRealm)

	body, _ := json.Marshal(newUser)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Failed to register user: %s", resp.Status)
	}
	return nil
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
