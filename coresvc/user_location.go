package coresvc

// UserLocation stores preferred location
type UserLocation struct {
	ID         uint `gorm:"primaryKey"`
	UserID     uint
	LocationID uint
}
