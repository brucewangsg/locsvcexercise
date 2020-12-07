package coresvc

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (r *routeContext) getLocationDetail(c *fiber.Ctx) error {
	location := &Location{}
	locationID, _ := strconv.Atoi(c.Params("id"))
	err := r.DB.Where("id = ?", locationID).First(location).Error

	if err != nil {
		c.SendStatus(http.StatusNotFound)
		return errors.New("no record found")
	}

	marshalledJSON, _ := json.Marshal(locationDetailJSONSerializer(*location))
	c.Send(marshalledJSON)
	return nil
}
