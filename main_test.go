package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brucewangsg/locsvcexercise/coresvc"
)

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

	user, err := createTestUser(db)
	if err != nil {
		t.Error("Failed to create user")
	}

	resp, _ := app.Test(httptest.NewRequest(
		"POST",
		"/auths/login",
		strings.NewReader(fmt.Sprintf(`{"email": "%s", "password": "PASS"}`, user.Email))),
	)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Login failed to return %v", http.StatusOK)
	}

	jsonResponseToken := getJwtToken(resp.Body)
	if jsonResponseToken == "" {
		t.Errorf("No token returned")
	}

	resp, _ = app.Test(httptest.NewRequest("POST", "/auths/verify", nil))
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Break in, it should require jwt token")
	}

	req := httptest.NewRequest("POST", "/auths/verify", nil)
	req.Header.Set("Authorization", "Bearer "+jsonResponseToken)
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

func TestLocationPreference(t *testing.T) {
	app, db := testApp()
	seedListingData(db)

	location := &coresvc.Location{}
	db.First(location)

	resp, _ := app.Test(httptest.NewRequest(
		"PUT",
		"/location_preference",
		strings.NewReader(fmt.Sprintf(`{"location_id": %d}`, location.ID))),
	)
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusOK, resp.StatusCode)
	}

	user, jwtToken := getUserAndJwtToken(app, db)
	req := httptest.NewRequest(
		"PUT",
		"/location_preference",
		strings.NewReader(fmt.Sprintf(`{"location_id": %d}`, location.ID)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Home page didn't return %v instead returning %v", http.StatusNoContent, resp.StatusCode)
	}

	updatedUserLocation := &coresvc.UserLocation{}
	db.Where("user_id = ?", user.ID).First(updatedUserLocation)
	if updatedUserLocation.LocationID != location.ID {
		t.Error("failed to update user preferred location")
	}
}
