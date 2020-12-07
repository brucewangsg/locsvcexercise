package authsvc

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type routeContext struct {
	DB *gorm.DB
}

// CurrentUser obtained from jwt token
type CurrentUser struct {
	Name  string
	Email string
	ID    uint
}

// JwtMiddleware intercept Authorization header and assign CurrentUser local
func JwtMiddleware(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		c.SendStatus(http.StatusUnauthorized)
		return errors.New("Unauthorized")
	}

	tokenString := ""
	if len(auth) > 7 {
		tokenString = auth[len("Bearer")+1:]
	}

	appSecret := os.Getenv("APP_SECRET")

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if token.Header["alg"] != "HS256" {
			return nil, fmt.Errorf("Invalid jwt token")
		}

		return []byte(appSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, _ := claims["id"].(float64)
		c.Locals("CurrentUser", &CurrentUser{
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
			ID:    uint(userID),
		})
		c.Next()
		return nil
	}

	c.SendStatus(http.StatusUnauthorized)
	return errors.New("Unauthorized")
}

// AddRoutes adding all endpoints to the fiber app
func AddRoutes(app *fiber.App, db *gorm.DB) {
	r := &routeContext{DB: db}

	app.Post("/auths/register", r.register)
	app.Post("/auths/login", r.login)

	app.Use("/auths/verify", JwtMiddleware)
	app.Get("/auths/verify", r.verify)
}
