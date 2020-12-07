package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brucewangsg/locsvcexercise/coresvc"
)

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
