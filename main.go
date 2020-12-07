package main

import (
	"fmt"
	"log"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/brucewangsg/locsvcexercise/coresvc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func newApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.SendString("Nothing to see here")
		return nil
	})
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${status} ${method} ${path}\n",
	}))
	coresvc.AddRoutes(app, db)
	authsvc.AddRoutes(app, db)

	return app
}

func main() {
	config := coresvc.NewAppConfig()
	db := coresvc.NewAppDBPool(config)

	app := newApp(db)
	log.Fatal(app.Listen(fmt.Sprintf(":%s", config.AppPort)))
}
