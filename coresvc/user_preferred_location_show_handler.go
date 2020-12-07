package coresvc

import (
	"encoding/json"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/gofiber/fiber/v2"
)

type userJSONSerializer authsvc.User

// MarshalJSON for userJSONSerializer
func (u userJSONSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  u.Name,
		Email: u.Email,
	})
}

func (r *routeContext) handleGetUserPreferredLocation(c *fiber.Ctx) error {
	currentUser := c.Locals("CurrentUser").(*authsvc.CurrentUser)

	userLocation := &UserLocation{}
	location := &Location{}
	r.DB.Where("user_id = ?", currentUser.ID).First(userLocation)
	r.DB.Where("id = ?", userLocation.LocationID).First(location)

	preferredLocationDetail := locationDetailJSONSerializer(*location)
	locationDetail := &preferredLocationDetail
	if location.ID == 0 {
		locationDetail = nil
	}

	userDetail := userJSONSerializer{
		Name:  currentUser.Name,
		Email: currentUser.Email,
	}
	marshalledJSON, _ := json.Marshal(&struct {
		User     userJSONSerializer            `json:"user"`
		Location *locationDetailJSONSerializer `json:"location,omitempty"`
	}{
		User:     userDetail,
		Location: locationDetail,
	})
	c.Send(marshalledJSON)

	return nil
}
