package coresvc

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type nextItemParams struct {
	LastBuildingName string `json:"last_building_name"`
	LastID           uint   `json:"last_id"`
}

func (r *routeContext) handleGetAllLocation(c *fiber.Ctx) error {
	params, _ := getListingParams(c)
	locations := remapToSerializedTypeLocations(
		*getFilteredLocations(r.DB, params),
	)

	marshalledJSON, _ := json.Marshal(locations)
	c.Send(marshalledJSON)
	return nil
}

func getFilteredLocations(db *gorm.DB, params *nextItemParams) *[]Location {
	locations := &[]Location{}
	finder := db.Limit(100)
	finder = finder.Order("building_name ASC, id DESC")

	if params.LastBuildingName != "" && params.LastID != 0 {
		finder = finder.Where(
			"building_name > ? OR (building_name = ? AND id < ?)",
			params.LastBuildingName,
			params.LastBuildingName,
			params.LastID)
	}
	finder.Find(locations)

	return locations
}

func getListingParams(c *fiber.Ctx) (*nextItemParams, error) {
	params := &nextItemParams{}
	reader := bytes.NewReader(c.Body())
	if err := json.NewDecoder(reader).Decode(params); err != nil {
		return params, err
	}
	return params, nil
}

func remapToSerializedTypeLocations(locations []Location) []locationJSONSerializer {
	remappedLocations := make([]locationJSONSerializer, 0, len(locations))
	for _, location := range locations {
		remappedLocations = append(remappedLocations, locationJSONSerializer(location))
	}
	return remappedLocations
}
