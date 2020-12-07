package authsvc

import (
	"encoding/json"
	"os"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
)

// User represents database fields to store user details
type User struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Email        string
	PasswordHash string
}

type userJSONSerializer User

func (u *User) getJwtToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = u.Name
	claims["email"] = u.Email
	claims["id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	appSecret := os.Getenv("APP_SECRET")

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}

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
