package models

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	GroupID   string `json:"groupId"`
	Function  string `json:"function"`
	Roles	[]string	`json:"roles"`

}

type KeycloakUser struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	FirstName string                 `json:"firstName"`
	LastName  string                 `json:"lastName"`
	Attributes map[string][]string   `json:"attributes"`
}

func (kc *KeycloakUser) ToUser() *User {
	return &User{
		ID:        kc.ID,
		Username:  kc.Username,
		FirstName: kc.FirstName,
		LastName:  kc.LastName,
		GroupID:   kc.getAttributeValue("groupId"),
		Function:  kc.getAttributeValue("function"),
	}
}

func (kc *KeycloakUser) getAttributeValue(key string) string {
	if values, exists := kc.Attributes[key]; exists && len(values) > 0 {
		return values[0]
	}
	return ""
}

type NewUser struct {
	Username               string `json:"username"`
	FirstName              string `json:"firstName"`
	LastName               string `json:"lastName"`
	Role                   string `json:"role"`
	Function               string `json:"function"`
	GroupID                string `json:"groupID"`
	Email                  string `json:"email"`
	Password               string `json:"password"`
	NotificationPreference bool   `json:"notificationPreference"`
}

func (n *NewUser) SetDefaultValues() {
	n.NotificationPreference = false
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"` //need to check on /role-mappings to see the naming
}

type KeycloakRole struct {
	ID string `json:"id"`
	Name string `json:"name"`
}