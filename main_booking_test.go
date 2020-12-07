package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brucewangsg/locsvcexercise/coresvc"
)

func TestBooking(t *testing.T) {
	app, db := testApp()
	seedListingData(db)

	location := &coresvc.Location{}
	db.Order("id DESC").First(location)
	location.AvailableSlot = 2
	db.Save(location)

	resp, _ := app.Test(httptest.NewRequest(
		"PUT",
		fmt.Sprintf("/bookings/%d", location.ID),
		nil))

	if resp.StatusCode == http.StatusOK {
		t.Error("booking for guest should not be allowed")
	}

	_, jwtToken := getUserAndJwtToken(app, db)
	req := httptest.NewRequest(
		"PUT",
		fmt.Sprintf("/bookings/%d", location.ID),
		nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	resp, _ = app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Error("booking failed")
	}

	db.Where("id = ?", location.ID).First(location)
	if location.AvailableSlot != 1 {
		t.Error("failed to reduce number of slot")
	}

	req = httptest.NewRequest(
		"PUT",
		fmt.Sprintf("/bookings/%d", location.ID),
		nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	resp, _ = app.Test(req)

	if resp.StatusCode != http.StatusOK {
		t.Error("it should restrict booking from the same person")
	}
}
