package coresvc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/gofiber/fiber/v2"
)

type updatePreferredLocationParams struct {
	LocationID uint `json:"location_id"`
}

func (r *routeContext) handleUpdateUserPreferredLocation(c *fiber.Ctx) error {
	params, _ := getUpdateParams(c)
	currentUser := c.Locals("CurrentUser").(*authsvc.CurrentUser)
	location := &Location{}

	if params.LocationID != 0 {
		err := r.DB.Where("id = ?", params.LocationID).First(location).Error
		if err != nil {
			c.SendStatus(404)
			return errors.New("location not found")
		}
	}

	if location.ID != 0 {
		preferredLocation := &UserLocation{}
		r.DB.FirstOrCreate(preferredLocation, UserLocation{UserID: currentUser.ID})
		preferredLocation.LocationID = location.ID
		r.DB.Save(preferredLocation)
	} else {
		c.SendStatus(404)
		return errors.New("location not found")
	}

	c.SendStatus(http.StatusNoContent)
	return nil
}

func getUpdateParams(c *fiber.Ctx) (*updatePreferredLocationParams, error) {
	params := &updatePreferredLocationParams{}
	reader := bytes.NewReader(c.Body())
	if err := json.NewDecoder(reader).Decode(params); err != nil {
		return params, err
	}
	return params, nil
}
