package coresvc

import (
	"github.com/brucewangsg/locsvcexercise/authsvc"
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
	app.Post("/locations", r.getAllLocation)
	app.Get("/locations/:id", r.getLocationDetail)

	app.Use("/location_preference", authsvc.JwtMiddleware)
	app.Put("/location_preference", r.updateUserPreferredLocation)
	app.Get("/location_preference", r.getUserPreferredLocation)
}
