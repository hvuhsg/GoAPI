package goapi

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestGetString(t *testing.T) {
	// Create a new request with a query parameter
	req, err := http.NewRequest("GET", "/?name=john", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]interface{}{},
	}

	// Add the query parameter to the Request object
	r.parameters["name"] = "john"

	// Test GetString with a valid parameter
	got, err := r.GetString("name")
	if err != nil {
		t.Errorf("GetString('name') returned an error: %v", err)
	}
	want := "john"
	if got != want {
		t.Errorf("GetString('name') = %q, want %q", got, want)
	}

	// Test GetString with an invalid parameter
	_, err = r.GetString("age")
	if err == nil {
		t.Errorf("GetString('age') did not return an error")
	}
}

func TestGetInt(t *testing.T) {
	// Create a new request with a query parameter
	req, err := http.NewRequest("GET", "/?age=25", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]interface{}{},
	}

	// Add the query parameter to the Request object
	r.parameters["age"] = 25

	// Test GetInt with a valid parameter
	got, err := r.GetInt("age")
	if err != nil {
		t.Errorf("GetInt('age') returned an error: %v", err)
	}
	want := 25
	if got != want {
		t.Errorf("GetInt('age') = %v, want %v", got, want)
	}

	// Test GetInt with an invalid parameter
	_, err = r.GetInt("name")
	if err == nil {
		t.Errorf("GetInt('name') did not return an error")
	}
}

func TestGetFloat(t *testing.T) {
	// Create a new request with a query parameter
	req, err := http.NewRequest("GET", "/?weight=75.5", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]interface{}{},
	}

	// Add the query parameter to the Request object
	r.parameters["weight"] = 75.5

	// Test GetFloat with a valid parameter
	got, err := r.GetFloat("weight")
	if err != nil {
		t.Errorf("GetFloat('weight') returned an error: %v", err)
	}
	want := 75.5
	if got != want {
		t.Errorf("GetFloat('weight') = %v, want %v", got, want)
	}

	// Test GetFloat with an invalid parameter
	_, err = r.GetFloat("name")
	if err == nil {
		t.Errorf("GetFloat('name') did not return an error")
	}
}

func TestGetBool(t *testing.T) {
	// Create a new request with a boolean parameter
	req, err := http.NewRequest("GET", "/?isAdult=true", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]interface{}{},
	}

	// Add the boolean parameter to the Request object
	r.parameters["isAdult"] = true

	// Test GetBool with a valid parameter
	got, err := r.GetBool("isAdult")
	if err != nil {
		t.Errorf("GetBool('isAdult') returned an error: %v", err)
	}
	want := true
	if got != want {
		t.Errorf("GetBool('isAdult') = %v, want %v", got, want)
	}

	// Test GetBool with an invalid parameter
	_, err = r.GetBool("age")
	if err == nil {
		t.Errorf("GetBool('age') did not return an error")
	}

	// Test GetBool with a non-boolean parameter
	r.parameters["age"] = 30
	_, err = r.GetBool("age")
	if err == nil {
		t.Errorf("GetBool('age') did not return an error")
	}
}

func TestGetArray(t *testing.T) {
	// Create a new request with an array parameter
	req, err := http.NewRequest("GET", "/?ids[]=1&ids[]=2&ids[]=3", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]interface{}{},
	}

	// Add the array parameter to the Request object
	r.parameters["ids"] = []interface{}{1, 2, 3}

	// Test GetArray with a valid parameter
	got, err := r.GetArray("ids")
	if err != nil {
		t.Errorf("GetArray('ids') returned an error: %v", err)
	}
	want := []interface{}{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetArray('ids') = %v, want %v", got, want)
	}

	// Test GetArray with an invalid parameter
	_, err = r.GetArray("names")
	if err == nil {
		t.Errorf("GetArray('names') did not return an error")
	}

	// Test GetArray with a non-array parameter
	r.parameters["names"] = "John, Jane, Jim"
	_, err = r.GetArray("names")
	if err == nil {
		t.Errorf("GetArray('names') did not return an error")
	}
}

func TestGetStringBoolMap(t *testing.T) {
	// Create a new request with a query parameter
	req, err := http.NewRequest("GET", "/?map=%7B%22foo%22%3Atrue%2C%22bar%22%3Afalse%7D", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a new Request object
	r := &Request{
		HTTPRequest: req,
		parameters:  map[string]any{},
	}

	// Add the query parameter to the Request object
	r.parameters["map"] = map[string]interface{}{
		"foo": true,
		"bar": false,
	}

	// Test GetStringBoolMap with a valid parameter
	got, err := r.GetStringBoolMap("map")
	if err != nil {
		t.Errorf("GetStringBoolMap('map') returned an error: %v", err)
	}
	want := map[string]bool{
		"foo": true,
		"bar": false,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetStringBoolMap('map') = %v, want %v", got, want)
	}

	// Test GetStringBoolMap with an invalid parameter
	_, err = r.GetStringBoolMap("foo")
	if err == nil {
		t.Errorf("GetStringBoolMap('foo') did not return an error")
	}
}

func TestNewRequest(t *testing.T) {
	// Create a request with query parameters and a JSON body
	queryParams := "name=john&age=30"
	jsonBody := `{"pets": {"dog": true, "cat": false}, "car": "audi"}`
	req, err := http.NewRequest("POST", "/test?"+queryParams, bytes.NewBuffer([]byte(jsonBody)))

	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Test the NewRequest function
	r, err := NewRequest(req)
	if err != nil {
		t.Fatalf("NewRequest returned an error: %v", err)
	}

	// Test query parameters
	name, err := r.GetString("name")
	if err != nil {
		t.Errorf("GetString('name') returned an error: %v", err)
	}
	wantName := "john"
	if name != wantName {
		t.Errorf("GetString('name') = %q, want %q", name, wantName)
	}

	age, err := r.GetInt("age")
	if err != nil {
		t.Errorf("GetInt('age') returned an error: %v", err)
	}
	wantAge := 30
	if age != wantAge {
		t.Errorf("GetInt('age') = %d, want %d", age, wantAge)
	}

	// Test body parameters
	car, err := r.GetString("car")
	if err != nil {
		t.Errorf("GetInt('age') returned an error: %v", err)
	}
	wantCar := "audi"
	if car != wantCar {
		t.Errorf("GetString('car') = %s, want %s", car, wantCar)
	}

	pets, err := r.GetStringBoolMap("pets")
	if err != nil {
		t.Errorf("GetStringBoolMap('pets') returned an error: %v", err)
	}
	wantPets := map[string]bool{
		"dog": true,
		"cat": false,
	}
	if !reflect.DeepEqual(pets, wantPets) {
		t.Errorf("GetStringBoolMap('pets') = %v, want %v", pets, wantPets)
	}
}
