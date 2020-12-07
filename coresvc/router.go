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

	app.Get("/locations", r.handleGetAllLocation)
	app.Post("/locations", r.handleGetAllLocation)
	app.Get("/locations/:id", r.handleGetLocationDetail)

	app.Use("/location_preference", authsvc.JwtMiddleware)
	app.Put("/location_preference", r.handleUpdateUserPreferredLocation)
	app.Get("/location_preference", r.handleGetUserPreferredLocation)
}
