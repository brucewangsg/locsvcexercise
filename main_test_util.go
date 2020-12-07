package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"

	"github.com/brucewangsg/locsvcexercise/authsvc"
	"github.com/brucewangsg/locsvcexercise/coresvc"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDBPool(config *coresvc.AppConfig) *gorm.DB {
	dbstr := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		config.DatabaseHost,
		config.DatabasePort,
		config.DatabaseUser,
		config.DatabasePass,
		config.DatabaseName+"_test",
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbstr,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic("Something is wrong with database")
	}

	return db
}

func testApp() (*fiber.App, *gorm.DB) {
	config := coresvc.NewAppConfig()
	db := newTestDBPool(config)

	db.Exec("TRUNCATE users")
	db.Exec("TRUNCATE locations")
	return newApp(db), db
}

func seedListingData(db *gorm.DB) {
	db.Exec(`
		INSERT INTO locations(building_name, address, city, country, phone_number) VALUES
			('Cyber Cafe X', 'Rose Blooming Town Street', 'Singapore', 'Singapore', '65432111'),
			('Big City Mall', 'Red Hill Road 03-44', 'Penang', 'Malaysia', '8387133'),
			('Toast Link Town', 'Crepe Seed Street 11-11', 'Jakarta', 'Indonesia', '99213911'),
			('Mighty House', 'Blue Street', 'Bangkok', 'Thailand', '78432111'),
			('Silent Cave', 'Green Street', 'Hanoi', 'Vietnam', '52332322'),
			('Cross Junction', 'Red Tower Street', 'Perth', 'Australia', '87432111');
	`)
}

func createTestUser(db *gorm.DB) (*authsvc.User, error) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("PASS"), 8)
	user := &authsvc.User{
		Name:         "User 1",
		Email:        "user@gmail.com",
		PasswordHash: string(hashedPass),
	}

	err := db.Create(user).Error
	return user, err
}

func getJwtToken(respBody io.ReadCloser) string {
	jsonResponse := &struct {
		Token string `json:"token"`
	}{}
	content, _ := ioutil.ReadAll(respBody)
	json.Unmarshal(content, jsonResponse)

	return jsonResponse.Token
}

func getUserAndJwtToken(app *fiber.App, db *gorm.DB) (*authsvc.User, string) {
	user, _ := createTestUser(db)
	resp, _ := app.Test(httptest.NewRequest(
		"POST",
		"/auths/login",
		strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "PASS"}`, user.Email))),
	)
	jsonResponseToken := getJwtToken(resp.Body)
	return user, jsonResponseToken
}
