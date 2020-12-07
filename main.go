package main

import (
	"fmt"

	"github.com/brucewangsg/locsvcexercise/coresvc"
)

func main() {
	fmt.Println("Hello App")
	config := coresvc.NewAppConfig()
	db := coresvc.NewAppDBPool(config)

	fmt.Println(db)
}
