package main

import (
	"fmt"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/brucewangsg/locsvcexercise/coresvc"
)

func main() {
	config := newAppConfig()
	db := newAppDBPool(config)

	fmt.Println("Migrate users table")
	db.AutoMigrate(&authsvc.User{})
	db.Exec("CREATE UNIQUE INDEX users_email ON users (email)")

	fmt.Println("Migrate locations table")
	db.AutoMigrate(&coresvc.Location{})
}
