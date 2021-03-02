package geolocation

import (
	"net/http"
	"reflect"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var response = Info{
	Status:         "success",
	Message:        "",
	Continent:      "North America",
	ContinentCode:  "NA",
	Country:        "United States",
	CountryCode:    "US",
	Region:         "CA",
	RegionName:     "California",
	City:           "San Francisco",
	District:       "",
	ZIP:            "94103",
	Latitude:       37.774929,
	Longitude:      -122.419418,
	Timezone:       "America/Los_Angeles",
	TimezoneOffset: -28800,
	ISP:            "Test Inc.",
	ORG:            "Test Inc.",
	ASN:            "AS0000 Test Inc.",
	Mobile:         false,
	Proxy:          false,
	IP:             "0.0.0.0",
}

func TestValidate(t *testing.T) {
	if err := response.ValidateFields("Continent"); err != nil {
		t.Errorf("Got unexpected error: %s", err)
	}

	if err := response.ValidateFields("Country", "All"); err != nil {
		t.Errorf("Got unexpected error: %s", err)
	}

	expected := "Invalid field name Status"
	if err := response.ValidateFields("Status"); err.Error() != expected {
		t.Errorf("Expected error %s, Got %s", expected, err)
	}

	expected = "Invalid field name DoesNotExist"
	if err := response.ValidateFields("DoesNotExist"); err.Error() != expected {
		t.Errorf("Expected error %s, Got %s", expected, err)
	}
}

func TestSetupRequest(t *testing.T) {
	request, err := setupRequest()
	if err != nil {
		t.Errorf("Got unexepected error: %s", err)
	}

	if request.Method != http.MethodGet {
		t.Errorf("Expected method: %s, Got %s", http.MethodGet, request.Method)
	}

	if request.URL.String() != "http://ip-api.com/json/?fields=37482495" {
		t.Errorf("Expected URL to be http://ip-api.com/json/?fields=37482495, Got: %s", request.URL)
	}

	if regexp.MustCompile(`^meridian/\d+\d+\d+$`).MatchString(request.Header.Get("User-Agent")) {
		t.Errorf("Expected User Agent to be in format meridian/0.0.0, Got %s", request.Header.Get("User-Agent"))
	}

	request, err = setupRequest("google.com")
	if err != nil {
		t.Errorf("Got unexepected error: %s", err)
	}

	if request.URL.String() != "http://ip-api.com/json/google.com?fields=37482495" {
		t.Errorf("Expected URL to be http://ip-api.com/json/google.com?fields=37482495, Got: %s", request.URL)
	}
}

func TestFieldToString(t *testing.T) {
	if v := fieldToString(reflect.ValueOf("test")); v != "test" {
		t.Errorf("String: Expected 'test', Got '%s'", v)
	}

	if v := fieldToString(reflect.ValueOf(1)); v != "1" {
		t.Errorf("Int: Expected '1', Got '%s'", v)
	}

	if v := fieldToString(reflect.ValueOf(3.141592)); v != "3.1416" {
		t.Errorf("Float: Expected '3.1416', Got '%s'", v)
	}

	if v := fieldToString(reflect.ValueOf(true)); v != "true" {
		t.Errorf("Bool: Expected 'true', Got '%s'", v)
	}

	if v := fieldToString(reflect.ValueOf([]int{1, 2})); v != "" {
		t.Errorf("Bool: Expected '', Got '%s'", v)
	}
}

func TestCanUseField(t *testing.T) {
	field, _ := reflect.Indirect(reflect.ValueOf(response)).Type().FieldByName("Country")

	// Test default true state.
	if !canUseField(field, map[string]bool{"Country": true}, false) {
		t.Error("Field 'Country' should be usable, but is marked not usable.")
	}

	// Test field already being marked.
	if canUseField(field, map[string]bool{"Country": false}, false) {
		t.Error("Field 'Country' should not be usable, but is marked as usable.")
	}

	// Test override in case field was marked.
	if !canUseField(field, map[string]bool{"Country": false}, true) {
		t.Error("Field 'Country' should be usable, but is marked not usable.")
	}

	field, _ = reflect.Indirect(reflect.ValueOf(response)).Type().FieldByName("Status")

	// Test disabled field.
	if canUseField(field, map[string]bool{"Status": true}, false) {
		t.Error("Field 'Status' should not be usable b/c it is disabled, but is marked as usable.")
	}

	// Test disabled field even if override is set to `true`.
	if canUseField(field, map[string]bool{"Status": true}, true) {
		t.Error("Field 'Status' should not be usable b/c it is disabled, but is marked as usable.")
	}
}

func TestStringFull(t *testing.T) {
	const responseStringFull = `Continent: North America
ContinentCode: NA
Country: United States
CountryCode: US
Region: CA
RegionName: California
City: San Francisco
District: 
ZIP: 94103
Latitude: 37.7749
Longitude: -122.4194
Timezone: America/Los_Angeles
TimezoneOffset: -28800
ISP: Test Inc.
ORG: Test Inc.
ASN: AS0000 Test Inc.
Mobile: false
Proxy: false
IP: 0.0.0.0
`

	if response.String() != responseStringFull {
		t.Error(cmp.Diff(responseStringFull, response.String()))
	}

	if response.ToString("All") != responseStringFull {
		t.Errorf(cmp.Diff(responseStringFull, response.ToString("All")))
	}
}

func TestSubsetString(t *testing.T) {
	const responseString = `Country: United States
RegionName: California
City: San Francisco
Latitude: 37.7749
Longitude: -122.4194
IP: 0.0.0.0
`

	cmd := response.ToString("Country", "RegionName", "City", "Latitude", "Longitude", "IP")
	if cmd != responseString {
		t.Error(cmp.Diff(responseString, cmd))
	}
}

func TestToJSON(t *testing.T) {
	var responseJSON = []byte(`{"country":"United States","zip":"94103","lat":0,"lon":0,"offset":0,"mobile":false,"proxy":false}`)

	v, err := response.ToJSON("Country", "ZIP", "Mobile")
	if err != nil {
		t.Errorf("Got unexepected error: %s", err)
	}

	if !cmp.Equal(responseJSON, v) {
		t.Error(cmp.Diff(string(responseJSON), string(v)))
	}
}
