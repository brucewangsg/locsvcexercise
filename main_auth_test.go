package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
