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
