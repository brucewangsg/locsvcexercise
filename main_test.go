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

	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM locations")
	return newApp(db), db
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
