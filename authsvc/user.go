package authsvc

// User represents database fields to store user details
type User struct {
	ID    uint `gorm:"primaryKey"`
	Name  string
	Email string
}
