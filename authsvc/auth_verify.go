package authsvc

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func (r *routeContext) verify(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(*CurrentUser)

	marshalledJSON, _ := json.Marshal(&userJSONSerializer{
		Name:  currentUser.Name,
		Email: currentUser.Email,
		ID:    currentUser.ID,
	})
	c.Send(marshalledJSON)

	return nil
}
