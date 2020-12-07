package authsvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *routeContext) handleLogin(c *fiber.Ctx) error {
	params, err := getLoginParams(c)

	user := &User{ID: 0}
	err = r.DB.Where("email = ?", params.Email).Find(user).Error
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password)); err != nil {
		c.SendStatus(http.StatusUnauthorized)
		return errors.New("Unathorized")
	}

	jwtToken, err := user.getJwtToken()
	if err != nil {
		return err
	}

	marshalledJSON, _ := json.Marshal(&struct {
		Token string
		User  userJSONSerializer
	}{
		Token: jwtToken,
		User:  userJSONSerializer(*user),
	})
	c.Send(marshalledJSON)

	return nil
}

func getLoginParams(c *fiber.Ctx) (*registerParams, error) {
	params := &registerParams{}
	reader := bytes.NewReader(c.Body())
	if err := json.NewDecoder(reader).Decode(params); err != nil {
		return params, err
	}

	return params, nil
}
