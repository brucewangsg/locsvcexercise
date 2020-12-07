package main

import (
	"fmt"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/brucewangsg/locsvcexercise/coresvc"
	"gorm.io/gorm"
)

func seedListingData(db *gorm.DB) {
	db.Exec(`
		INSERT INTO locations(building_name, address, city, country, phone_number, available_slot) VALUES
			('Cyber Cafe X', 'Rose Blooming Town Street', 'Singapore', 'Singapore', '65432111', 100),
			('Big City Mall', 'Red Hill Road 03-44', 'Penang', 'Malaysia', '8387133', 100),
			('Toast Link Town', 'Crepe Seed Street 11-11', 'Jakarta', 'Indonesia', '99213911', 100),
			('Mighty House', 'Blue Street', 'Bangkok', 'Thailand', '78432111', 100),
			('Silent Cave', 'Green Street', 'Hanoi', 'Vietnam', '52332322', 100),
			('Cross Junction', 'Red Tower Street', 'Perth', 'Australia', '87432111', 100);
	`)
}

func main() {
	config := coresvc.NewAppConfig()
	db := coresvc.NewAppDBPool(config)

	fmt.Println("Migrate users table")
	db.AutoMigrate(&authsvc.User{})
	db.Exec("CREATE UNIQUE INDEX users_email ON users (email)")

	fmt.Println("Migrate locations table")
	db.AutoMigrate(&coresvc.Location{})
	db.Exec("CREATE INDEX location_names ON users (building_name)")
	seedListingData(db)

	fmt.Println("Migrate user locations table")
	db.Exec("CREATE UNIQUE INDEX uniq_preferred_locations ON users (user_id)")
	db.AutoMigrate(&coresvc.UserLocation{})

	fmt.Println("Migrate bookings table")
	db.AutoMigrate(&coresvc.Booking{})
	db.Exec("CREATE UNIQUE INDEX uniq_bookings ON bookings (user_id, location_id)")
}
