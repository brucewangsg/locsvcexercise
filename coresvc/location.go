package coresvc

import "encoding/json"

// Location represents database fields to store location details
type Location struct {
	ID            uint `gorm:"primaryKey"`
	BuildingName  string
	Address       string
	City          string
	Country       string
	PhoneNumber   string
	AvailableSlot int
}

type locationJSONSerializer Location
type locationDetailJSONSerializer Location

// MarshalJSON for locationJSONSerializer
func (l locationJSONSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint   `json:"id"`
		BuildingName string `json:"name"`
	}{
		ID:           l.ID,
		BuildingName: l.BuildingName,
	})
}

// MarshalJSON for locationDetailJSONSerializer
func (l locationDetailJSONSerializer) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID           uint   `json:"id"`
		BuildingName string `json:"name"`
		Address      string `json:"address"`
		City         string `json:"city"`
		Country      string `json:"country"`
		PhoneNumber  string `json:"phone_number"`
	}{
		ID:           l.ID,
		BuildingName: l.BuildingName,
		Address:      l.Address,
		City:         l.City,
		Country:      l.Country,
		PhoneNumber:  l.PhoneNumber,
	})
}
