package main

import (
	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/brucewangsg/locsvcexercise/coresvc"
)

func main() {
	config := newAppConfig()
	db := newAppDBPool(config)

	db.AutoMigrate(&authsvc.User{})
	db.Exec("CREATE UNIQUE INDEX users_email ON users (email)")

	db.AutoMigrate(&coresvc.Location{})
}
