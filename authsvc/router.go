package authsvc

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type routeContext struct {
	DB *gorm.DB
}

// AddRoutes adding all endpoints to the fiber app
func AddRoutes(app *fiber.App, db *gorm.DB) {
	r := &routeContext{DB: db}

	app.Post("/auths/register", r.register)
	app.Post("/auths/login", r.login)
	app.Post("/auths/verify", r.verify)
}
