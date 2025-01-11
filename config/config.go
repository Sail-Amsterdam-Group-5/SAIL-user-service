package config

type Config struct {
	KeycloakURL   string
	KeycloakRealm string
	AdminToken    string
	Port          string
}

func LoadConfig() (*Config, error) {
	return &Config{
		KeycloakURL:   "http://keycloak-route-oscar-dev.apps.inholland.hcs-lab.nl",
		KeycloakRealm: "sail-amsterdam",
		AdminToken:    "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJmNDJyd1plbXRyWUdabWYzQjE1eklvZlhzQ3FnV0hDUUloMHBDWERURWNVIn0.eyJleHAiOjE3MzY2NDA1NTUsImlhdCI6MTczNjYwNDU1NSwianRpIjoiNjk3ZDRjMWMtYzc0My00YzVhLTliYTUtZTQ3OTVlN2VmMDlmIiwiaXNzIjoiaHR0cDovL2tleWNsb2FrLXJvdXRlLW9zY2FyLWRldi5hcHBzLmluaG9sbGFuZC5oY3MtbGFiLm5sL3JlYWxtcy9zYWlsLWFtc3RlcmRhbSIsImF1ZCI6WyJyZWFsbS1tYW5hZ2VtZW50IiwiYWNjb3VudCJdLCJzdWIiOiIzYjZmYTVlZS0zZmJjLTRmZTQtYjFmOS0wMWQ1OGUzODJmMzQiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ1c2VyLW1pY3Jvc2VydmljZSIsImFjciI6IjEiLCJhbGxvd2VkLW9yaWdpbnMiOlsiLyoiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwidW1hX2F1dGhvcml6YXRpb24iLCJkZWZhdWx0LXJvbGVzLXNhaWwtYW1zdGVyZGFtIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsicmVhbG0tbWFuYWdlbWVudCI6eyJyb2xlcyI6WyJtYW5hZ2UtdXNlcnMiLCJ2aWV3LXVzZXJzIiwicXVlcnktZ3JvdXBzIiwicXVlcnktdXNlcnMiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoicHJvZmlsZSBlbWFpbCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwiY2xpZW50SG9zdCI6IjEwLjEyOC4yLjEyIiwicHJlZmVycmVkX3VzZXJuYW1lIjoic2VydmljZS1hY2NvdW50LXVzZXItbWljcm9zZXJ2aWNlIiwiY2xpZW50QWRkcmVzcyI6IjEwLjEyOC4yLjEyIiwiY2xpZW50X2lkIjoidXNlci1taWNyb3NlcnZpY2UifQ.lz-Zi7Qm33CwGz1QQ5A-aV1hvojrVjDMy7RbBDfO3Kll4PjXZ0boPkflhEIItKUGRPRhZsgmtPxE-SfT0TtSeVeRXPuXchRAIWZwUJ8bFY1z3upqWvqbgrX5b-2RfoGaJbaUs_TLDlVIyeFuPFkKXSIMm6or45lhfp-tUZvfvpLfgJq9H-42J3yWgQE52ocSKPcX4HQU8Nv68vuDeJNxJJlcpfQpNt8-6EAm6VPPzeSUz7wZ1HejF0Lv-bMUF2BFkKiMzDPOZKG9lNIfow8_n8tVY9iBq8bWhKIwb_CMw7TdZy_m5nhi8EEmNO788-D4dCt-_B7ve59rZAiVHTmR-g",
		Port:          "8080",
	}, nil
}
