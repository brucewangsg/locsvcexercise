package authsvc

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type registerParams struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *routeContext) register(c *fiber.Ctx) error {
	params, err := getRegisterParams(c)
	if err != nil {
		return err
	}

	user := &User{
		Name:  params.Name,
		Email: params.Email,
	}

	err = r.DB.Create(user).Error
	if err != nil {
		return err
	}

	c.SendStatus(201)

	return nil
}

func getRegisterParams(c *fiber.Ctx) (*registerParams, error) {
	params := &registerParams{}
	reader := bytes.NewReader(c.Body())
	if err := json.NewDecoder(reader).Decode(params); err != nil {
		return params, err
	}
	if params.Password == "" {
		return params, errors.New("missing password for new user")
	}
	return params, nil
}
