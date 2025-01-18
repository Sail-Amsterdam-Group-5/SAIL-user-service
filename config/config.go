package config

import "os"

type Config struct {
	KeycloakURL   string
	KeycloakRealm string
	Port          string
	ClientID      string
	ClientSecret  string
}

func LoadConfig() (*Config, error) {
	return &Config{
		KeycloakURL:   os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm: "sail-amsterdam",
		Port:          "8080",
		ClientID:      os.Getenv("CLIENT_ID"),
		ClientSecret:  os.Getenv("CLIENT_SECRET"),
	}, nil
}
