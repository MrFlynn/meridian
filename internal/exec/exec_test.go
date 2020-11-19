package exec

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mrflynn/meridian/geolocation"
)

var response = &geolocation.Info{
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

func TestParseCommandString(t *testing.T) {
	cmdStr, err := ParseCommandString("echo {{ .Country }} {{ .City }}", response)
	if err != nil {
		t.Errorf("Got unexepected error: %s", err)
	}

	expected := []string{"echo", "United States", "San Francisco"}
	if cmp.Equal(expected, cmdStr) {
		t.Error(cmp.Diff(expected, cmdStr))
	}

	_, err = ParseCommandString("echo {{ .Test }} {{ .ISP }}", response)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
