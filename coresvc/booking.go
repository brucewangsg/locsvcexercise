package coresvc

// Booking stores location bookings
type Booking struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	LocationID uint
}
