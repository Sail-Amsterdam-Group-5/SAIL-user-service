package config

type Config struct {
	KeycloakURL   string
	KeycloakRealm string
	Port          string
}

func LoadConfig() (*Config, error) {
	return &Config{
		KeycloakURL:   "http://keycloak-route-oscar-dev.apps.inholland.hcs-lab.nl", //GET RID OF HARDCODED URL
		KeycloakRealm: "sail-amsterdam",
		Port:          "8080",
	}, nil
}
