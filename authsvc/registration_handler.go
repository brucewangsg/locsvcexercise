package authsvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type registerParams struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r *routeContext) handleRegister(c *fiber.Ctx) error {
	params, err := getRegisterParams(c)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)
	if err != nil {
		return err
	}

	user := &User{
		Name:         params.Name,
		Email:        params.Email,
		PasswordHash: string(hashedPassword),
	}

	err = r.DB.Create(user).Error
	if err != nil {
		return err
	}

	c.SendStatus(201)

	return nil
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func getRegisterParams(c *fiber.Ctx) (*registerParams, error) {
	params := &registerParams{}
	reader := bytes.NewReader(c.Body())
	if err := json.NewDecoder(reader).Decode(params); err != nil {
		return params, err
	}

	if params.Password == "" {
		return params, errors.New("missing password for new user")
	}

	if params.Name == "" {
		return params, errors.New("missing name for new user")
	}

	if len(params.Name) >= 30 {
		return params, errors.New("name is too long for new user")
	}

	if params.Email == "" {
		return params, errors.New("missing email for new user")
	}

	if len(params.Email) < 3 ||
		len(params.Email) >= 60 ||
		!emailRegex.MatchString(params.Email) {
		return params, errors.New("invalid email format for new user")
	}

	return params, nil
}
