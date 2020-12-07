package coresvc

// Location represents database fields to store location details
type Location struct {
	ID           uint `gorm:"primaryKey"`
	BuildingName string
	Address      string
	City         string
	Country      string
	PhoneNumber  string
}
