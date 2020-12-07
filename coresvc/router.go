package coresvc

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

	app.Get("/locations", r.getAllLocation)
	app.Get("/locations/:id", r.getLocationDetail)
}

func (r *routeContext) getAllLocation(c *fiber.Ctx) error {
	c.SendString("{\"info\": \"list down all locations\"}")
	return nil
}

func (r *routeContext) getLocationDetail(c *fiber.Ctx) error {
	c.SendString("{\"info\": \"get location details\"}")
	return nil
}
