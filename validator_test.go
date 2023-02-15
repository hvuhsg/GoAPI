package goapi

import (
	"testing"
)

func TestVRequired(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": 42}}
	validator := VRequired{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsString(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": "hello"}}
	validator := VIsString{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = 42
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsInt(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": 42}}
	validator := VIsInt{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = "not an int"
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsFloat(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": 3.14}}
	validator := VIsFloat{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = "not a float"
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsBool(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": true}}
	validator := VIsBool{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = "not a bool"
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsArray(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": []interface{}{1, 2, 3}}}
	validator := VIsArray{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = "not an array"
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestVIsMap(t *testing.T) {
	req := &Request{parameters: map[string]interface{}{"param1": map[string]interface{}{"a": 1, "b": 2}}}
	validator := VIsMap{}

	err := validator.Validate(req, "param1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	req.parameters["param2"] = "not a map"
	err = validator.Validate(req, "param2")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
