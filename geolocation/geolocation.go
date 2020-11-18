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

	request.Header.Set("User-Agent", "meridian/0.1.0")

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
	Continent      string  `json:"continent,omitempty"`
	ContinentCode  string  `json:"continentCode,omitempty"`
	Country        string  `json:"country,omitempty"`
	CountryCode    string  `json:"countryCode,omitempty"`
	Region         string  `json:"region,omitempty"`
	RegionName     string  `json:"regionName,omitempty"`
	City           string  `json:"city,omitempty"`
	District       string  `json:"district,omitempty"`
	ZIP            string  `json:"zip,omitempty"`
	Latitude       float64 `json:"lat,omitempty"`
	Longitude      float64 `json:"lon,omitempty"`
	Timezone       string  `json:"timezone,omitempty"`
	TimezoneOffset int     `json:"offset,omitempty"`
	ISP            string  `json:"isp,omitempty"`
	ORG            string  `json:"org,omitempty"`
	ASN            string  `json:"as,omitempty"`
	Mobile         bool    `json:"mobile,omitempty"`
	Proxy          bool    `json:"proxy,omitempty"`
	IP             string  `json:"query,omitempty"`
}

// New takes an optional location (IP address or domain name) and returns
// geolocation information about the location.
func New(location ...string) (*Info, error) {
	request, err := setupRequest(location...)
	if err != nil {
		return nil, err
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var info Info
	err = json.Unmarshal(content, &info)
	if err != nil {
		return nil, err
	}

	if info.Status != "success" {
		return nil, errors.New(info.Message)
	}

	return &info, nil
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
	aux := Info{}
	auxRef := reflect.ValueOf(aux)

	it := i.toIter(fields...)
	for it.next() {
		value, field, err := it.value()
		if err != nil {
			continue
		}

		if targetField := auxRef.FieldByName(field.Name); targetField.CanSet() {
			targetField.Set(value)
		}
	}

	return json.Marshal(aux)
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
