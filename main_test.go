package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestApp(t *testing.T) {
	app, _ := testApp()

	resp, _ := app.Test(httptest.NewRequest("GET", "/", nil))
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestRegistration(t *testing.T) {
	app, _ := testApp()

	resp, _ := app.Test(httptest.NewRequest(
		"POST",
		"/auths/register",
		strings.NewReader(`{"email": "user@gmail.com", "name": "User 1", "password": "PASS"}`)),
	)
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusCreated, resp.StatusCode)
	}

	resp, _ = app.Test(httptest.NewRequest(
		"POST",
		"/auths/register",
		strings.NewReader(`{"email": "user@gmail.com", "name": "User 1", "password": "PASS"}`)),
	)

	content, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(content), "duplicate key value") {
		t.Errorf("should not create due to duplication constraint, %v", content)
	}
}

func TestLogin(t *testing.T) {
	app, db := testApp()

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("PASS"), 8)
	user := &authsvc.User{
		Name:         "User 1",
		Email:        "user@gmail.com",
		PasswordHash: string(hashedPass),
	}

	err := db.Create(user).Error
	if err != nil {
		t.Error("Failed to create user")
	}

	resp, _ := app.Test(httptest.NewRequest(
		"POST",
		"/auths/login",
		strings.NewReader(`{"email": "user@gmail.com", "password": "PASS"}`)),
	)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Login failed to return %v", http.StatusOK)
	}

	jsonResponse := &struct {
		Token string `json:"token"`
	}{}
	content, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(content, jsonResponse)

	if jsonResponse.Token == "" {
		t.Errorf("No token returned")
	}

	resp, _ = app.Test(httptest.NewRequest("POST", "/auths/verify", nil))
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Break in, it should require jwt token")
	}

	req := httptest.NewRequest("POST", "/auths/verify", nil)
	req.Header.Set("Authorization", "Bearer "+jsonResponse.Token)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("It should allow login with supplied jwt token")
	}

}

func TestLocationListing(t *testing.T) {
	app, db := testApp()
	seedListingData(db)

	resp, _ := app.Test(httptest.NewRequest(
		"GET",
		"/locations",
		nil),
	)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusOK, resp.StatusCode)
	}

	content, _ := ioutil.ReadAll(resp.Body)
	if string(content) == "[]" {
		t.Error("it should return some locations")
	}

	location := &coresvc.Location{}
	db.Order("id DESC").First(location)

	resp, _ = app.Test(httptest.NewRequest(
		"POST",
		"/locations",
		strings.NewReader(fmt.Sprintf(`{"last_building_name": "XTown", "last_id": %d}`, location.ID))),
	)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusOK, resp.StatusCode)
	}

	content, _ = ioutil.ReadAll(resp.Body)
	if string(content) != "[]" {
		t.Errorf("it should return empty list instead of %v", string(content))
	}
}

func TestLocationDetail(t *testing.T) {
	app, db := testApp()
	seedListingData(db)

	location := &coresvc.Location{}
	db.First(location)

	resp, _ := app.Test(httptest.NewRequest(
		"GET",
		fmt.Sprintf("/locations/%d", location.ID),
		nil),
	)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusOK, resp.StatusCode)
	}
}
