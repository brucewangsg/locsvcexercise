package authsvc

import "encoding/json"

// User represents database fields to store user details
type User struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Email        string
	PasswordHash string
}

type userJSONSerializer User

// MarshalJSON for userJSONSerializer
func (u userJSONSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  u.Name,
		Email: u.Email,
	})
}
