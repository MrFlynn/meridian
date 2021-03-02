package geolocation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

const endpointTemplate = "http://ip-api.com/json/%s?fields=37482495"

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

func setupRequest(location ...string) (*http.Request, error) {
	var loc string
	if len(location) > 0 {
		loc = location[0]
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf(endpointTemplate, loc), nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("User-Agent", fmt.Sprintf("meridian/%s", viper.GetString("version")))

	return request, nil
}

func fieldToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Float64, reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', 4, 64)
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	default:
		return ""
	}
}

func canUseField(field reflect.StructField, useFields map[string]bool, skipCheck bool) bool {
	status, ok := useFields[field.Name]
	return ((status && ok) || skipCheck) && field.Tag.Get("meridian") != "disable"
}

// Info stores location query results.
type Info struct {
	Status         string  `json:"status,omitempty" meridian:"disable"`
	Message        string  `json:"message,omitempty" meridian:"disable"`
	Continent      string  `json:"continent,omitempty" description:"Full name of continent" ex:"North America"`
	ContinentCode  string  `json:"continentCode,omitempty" description:"Shorthand name of continent" ex:"NA"`
	Country        string  `json:"country,omitempty" description:"Full name of country" ex:"United States"`
	CountryCode    string  `json:"countryCode,omitempty" description:"Shorthand name of country" ex:"US"`
	Region         string  `json:"region,omitempty" description:"Shorthand name of region, state, etc." ex:"CA"`
	RegionName     string  `json:"regionName,omitempty" description:"Full name of region, state, etc." ex:"California"`
	City           string  `json:"city,omitempty" description:"Full name of city" ex:"San Francisco"`
	District       string  `json:"district,omitempty" description:"Full name of city district" ex:"South of Market"`
	ZIP            string  `json:"zip,omitempty" description:"Postal code" ex:"94103"`
	Latitude       float64 `json:"lat"`
	Longitude      float64 `json:"lon"`
	Timezone       string  `json:"timezone,omitempty" description:"tzdata name of timezone" ex:"America/Los_Angeles"`
	TimezoneOffset int     `json:"offset" description:"Offset in seconds from UTC" ex:"-28800 for America/Los_Angeles"`
	ISP            string  `json:"isp,omitempty" description:"Name of ISP" ex:"Comcast Cable Communications, LLC"`
	ORG            string  `json:"org,omitempty" description:"Organizational owner of IP, usually ISP" ex:"Comcast Cable Communications, Inc"`
	ASN            string  `json:"as,omitempty" description:"AS name and number for current IP" ex:"AS7922 Comcast Cable Communications, LLC"`
	Mobile         bool    `json:"mobile" description:"Whether or not you are on a mobile network"`
	Proxy          bool    `json:"proxy" description:"Whether or not you are using a proxy"`
	IP             string  `json:"query,omitempty" description:"Current IP address" ex:"0.0.0.0"`
}

// New takes an optional location (IP address or domain name) and returns
// geolocation information about the location.
func (i *Info) New(location ...string) error {
	request, err := setupRequest(location...)
	if err != nil {
		return err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return errors.New("Too many requests. Please wait 1 minute")
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Error: %s", response.Status)
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, i)
	if err != nil {
		return err
	}

	if i.Status != "success" {
		return errors.New(i.Message)
	}

	return nil
}

// ValidateFields validates that the given field names exist within the given Info object.
func (i *Info) ValidateFields(fields ...string) error {
	t := reflect.Indirect(reflect.ValueOf(i)).Type()

	for _, field := range fields {
		// If we encounter this special variable, we can skip further checks.
		if field == "All" {
			return nil
		}

		if sf, ok := t.FieldByName(field); !ok || sf.Tag.Get("meridian") == "disable" {
			return fmt.Errorf("Invalid field name %s", field)
		}
	}

	return nil
}

// ToString returns a formatted string representation of the Info struct.
// The optional fields argument takes a list of fields to output. Giving
// no arguments will output an empty string. Specifying "All" will return
// all fields.
func (i *Info) ToString(fields ...string) string {
	builder := strings.Builder{}

	it := i.toIter(fields...)
	for it.next() {
		value, field, err := it.value()
		if err != nil {
			continue
		}

		builder.WriteString(field.Name)
		builder.WriteString(": ")
		builder.WriteString(fieldToString(value))
		builder.WriteByte('\n')
	}

	return builder.String()
}

func (i *Info) String() string {
	return i.ToString("All")
}

// ToJSON takes a list of fields to convert to JSON byte array from a Info struct and returns
// the byte array with only those fields (unless the All specifier is given), regardless of
// whether or not the fields are set in the Info struct.
func (i *Info) ToJSON(fields ...string) ([]byte, error) {
	aux := &Info{}
	auxRef := reflect.ValueOf(aux)

	it := i.toIter(fields...)
	for it.next() {
		value, field, err := it.value()
		if err != nil {
			continue
		}

		if targetField := auxRef.Elem().FieldByName(field.Name); targetField.CanSet() {
			targetField.Set(value)
		}
	}

	return json.Marshal(&aux)
}

// ToDescription returns a formatted description of all fields in the Info struct.
func (i *Info) ToDescription() string {
	var (
		builder = &strings.Builder{}

		bold    = color.New(color.Bold)
		italics = color.New(color.Italic)
	)

	it := i.toIter("All")
	for it.next() {
		_, field, err := it.value()
		if err != nil {
			continue
		}

		builder.WriteString("- ")
		bold.Fprint(builder, field.Name)

		if description, ok := field.Tag.Lookup("description"); ok {
			builder.WriteString(": ")
			builder.WriteString(description)
		}

		if example, ok := field.Tag.Lookup("ex"); ok {
			builder.WriteString(", ")
			italics.Fprint(builder, "ex. ", example)
		}

		builder.WriteString(".\n")
	}

	return builder.String()
}

type infoIter struct {
	ref            reflect.Value
	idx, stop      int
	emitFields     map[string]bool
	skipFieldCheck bool
}

func (i *Info) toIter(fields ...string) *infoIter {
	value := reflect.Indirect(reflect.ValueOf(i))

	iter := &infoIter{
		ref:        value,
		idx:        0,
		stop:       value.NumField(),
		emitFields: make(map[string]bool, len(fields)),
	}

	for _, field := range fields {
		iter.emitFields[field] = true
	}

	_, iter.skipFieldCheck = iter.emitFields["All"]

	return iter
}

func (it *infoIter) next() bool {
	return it.idx < it.stop
}

func (it *infoIter) value() (reflect.Value, reflect.StructField, error) {
	defer func() {
		it.idx++
	}()

	if !it.next() {
		return reflect.Value{}, reflect.StructField{}, errors.New("iterator exausted")
	}

	value := it.ref.Field(it.idx)
	field := it.ref.Type().Field(it.idx)

	if !canUseField(field, it.emitFields, it.skipFieldCheck) {
		return reflect.Value{}, reflect.StructField{}, fmt.Errorf("field %s skipped", field.Name)
	}

	it.emitFields[field.Name] = false
	return value, field, nil
}
