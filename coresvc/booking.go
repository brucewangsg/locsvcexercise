package coresvc

// UserLocation stores preferred location
type Booking struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	LocationID uint
}
