package keycloak

import (
	"SAIL-user-service/config"
	"SAIL-user-service/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"log"

	"github.com/sirupsen/logrus"
)

type Client struct {
	config *config.Config
	logger *logrus.Entry
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func NewKeycloakClient(cfg *config.Config) *Client {
	logger := logrus.WithFields(logrus.Fields{
		"module": "keycloak",
	})
	return &Client{config: cfg, logger: logger}
}

func (kc *Client) GetClientToken() (string, error) {
	kc.logger.Info("Fetching client token")

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.config.KeycloakURL, kc.config.KeycloakRealm)

	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials",
		kc.config.ClientID, kc.config.ClientSecret)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		kc.logger.WithError(err).Error("Failed to create token request")
		return "", fmt.Errorf("internal server error")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		kc.logger.WithError(err).Error("Failed to fetch token")
		return "", fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		kc.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"url":         url,
		}).Error("Failed to fetch token")
		return "", fmt.Errorf("failed to fetch token: %s", resp.Status)
	}
	log.Println("Got token")

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		kc.logger.WithError(err).Error("Failed to decode token response")
		return "", fmt.Errorf("internal server error")
	}
	log.Println("Succesfully decoded token")

	kc.logger.Info("Successfully retrieved client token")
	return tokenResponse.AccessToken, nil
}

func (kc *Client) GetAllUsers() ([]models.User, error) {
	kc.logger.Info("Fetching all users")

	token, err := kc.GetClientToken()
	if err != nil {
		kc.logger.WithError(err).Error("Failed to retrieve client token")
		return nil, fmt.Errorf("internal server error")
	}
	log.Println("Got token")


	url := fmt.Sprintf("%s/admin/realms/%s/users", kc.config.KeycloakURL, kc.config.KeycloakRealm)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		kc.logger.WithError(err).Error("Failed to create request for all users")
		return nil, fmt.Errorf("internal server error")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	log.Println("Created request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		kc.logger.WithError(err).Error("Error while fetching users from Keycloak")
		return nil, fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		kc.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"url":         url,
		}).Error("Failed to fetch users")
		return nil, fmt.Errorf("failed to fetch users: %s", resp.Status)
	}

	var keycloakUsers []models.KeycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUsers); err != nil {
		kc.logger.WithError(err).Error("Failed to decode user list")
		return nil, fmt.Errorf("internal server error")
	}

	kc.logger.WithField("user_count", len(keycloakUsers)).Info("Successfully fetched all users")
	users := make([]models.User, 0, len(keycloakUsers))
	for _, keycloakUser := range keycloakUsers {
		//get user roles
		roles, err := kc.GetUserRoles(keycloakUser.ID)
		if err != nil {
			kc.logger.WithError(err).WithField("user_id", keycloakUser.ID).Warn("Failed to fetch roles for user")
			return nil, fmt.Errorf("internal server error")
		}

		user := keycloakUser.ToUser()
		user.Roles = roles
		users = append(users, *user)
	}

	return users, nil
}

func (kc *Client) GetUserById(id string) (*models.User, error) {
	kc.logger.WithField("user_id", id).Info("Fetching user by ID")

	token, err := kc.GetClientToken()
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to retrieve client token")
		return nil, fmt.Errorf("internal server error")
	}
	log.Println("Got token")

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.config.KeycloakURL, kc.config.KeycloakRealm, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to create request for user by ID")
		return nil, fmt.Errorf("internal server error")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	log.Println("Created request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to fetch user from Keycloak")
		return nil, fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		kc.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"user_id":     id,
			"url":         url,
		}).Error("Failed to fetch user")
		return nil, fmt.Errorf("failed to fetch user: %s", resp.Status)
	}

	var keycloakUser models.KeycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUser); err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to decode user data")
		return nil, fmt.Errorf("internal server error")
	}

	//adding the roles
	roles, err := kc.GetUserRoles(id)
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Warn("Failed to fetch roles for user")
		return nil, fmt.Errorf("internal server error")
	}

	user := keycloakUser.ToUser()
	user.Roles = roles

	log.Println("Returning user")
	
	return user, nil
}

// adding roles to the json bc the /users endpoint doesn't by default
// do this in O(1) lmao
func (kc *Client) GetUserRoles(id string) ([]string, error) {
	kc.logger.WithField("user_id", id).Info("Fetching roles for user")

	token, err := kc.GetClientToken()
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to retrieve client token")
		return nil, fmt.Errorf("internal server error")
	}

	url := fmt.Sprintf("%s/admin/realms/%s/users/%s/role-mappings/realm", kc.config.KeycloakURL, kc.config.KeycloakRealm, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to create request for user roles")
		return nil, fmt.Errorf("internal server error")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to fetch roles from Keycloak")
		return nil, fmt.Errorf("internal server error")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		kc.logger.WithFields(logrus.Fields{
			"status_code": resp.StatusCode,
			"user_id":     id,
			"url":         url,
		}).Error("Failed to fetch roles")
		return nil, fmt.Errorf("failed to fetch roles: %s", resp.Status)
	}

	var keycloakRoles []models.KeycloakRole
	if err := json.NewDecoder(resp.Body).Decode(&keycloakRoles); err != nil {
		kc.logger.WithError(err).WithField("user_id", id).Error("Failed to decode roles data")
		return nil, fmt.Errorf("internal server error")
	}

	roles := make([]string, len(keycloakRoles))
	for i, role := range keycloakRoles {
		roles[i] = role.Name
	}

	kc.logger.WithFields(logrus.Fields{
		"user_id":     id,
		"role_count":  len(roles),
	}).Info("Successfully fetched user roles")
	
	return roles, nil
}
