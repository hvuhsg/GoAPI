package request

import (
	"fmt"
	"reflect"
	"testing"
)

var panicValue interface{}

// Helper function to catch panics and return the panic value
func catchPanicString(f func() string) string {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered panic: %v", r)
			panicValue = err
		}
	}()
	panicValue = nil
	returnValue := f()
	if panicValue != nil {
		return ""
	}
	return returnValue
}

func catchPanicInt(f func() int) int {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered panic: %v", r)
			panicValue = err
		}
	}()
	panicValue = nil
	returnValue := f()
	if panicValue != nil {
		return 0
	}
	return returnValue
}

func catchPanicFloat(f func() float64) float64 {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered panic: %v", r)
			panicValue = err
		}
	}()
	panicValue = nil
	returnValue := f()
	if panicValue != nil {
		return 0
	}
	return returnValue
}

func catchPanicBool(f func() bool) bool {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered panic: %v", r)
			panicValue = err
		}
	}()
	panicValue = nil
	returnValue := f()
	if panicValue != nil {
		return false
	}
	return returnValue
}

func catchPanicArray(f func() []int) []int {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered panic: %v", r)
			panicValue = err
		}
	}()
	panicValue = nil
	returnValue := f()
	if panicValue != nil {
		return nil
	}
	return returnValue
}

func TestGetString(t *testing.T) {
	t.Cleanup(func() {
		panicValue = nil
	})

	// Define a test request with some parameters
	req := &Request{
		Parameters: map[string]interface{}{
			"str": "hello",
			"num": 42,
		},
	}

	// Test case 1: get an existing string parameter
	strVal := catchPanicString(func() string { return req.GetString("str") })
	if panicValue != nil {
		t.Errorf("GetString panic'd: %v", panicValue)
	} else if strVal != "hello" {
		t.Errorf("GetString returned wrong value: expected 'hello', got '%s'", strVal)
	}

	// Test case 2: get a non-existing parameter
	nonExistingVal := catchPanicString(func() string { return req.GetString("non-existing") })
	if panicValue == nil {
		t.Errorf("GetString did not panic for non-existing parameter")
	} else if nonExistingVal != "" {
		t.Errorf("GetString returned wrong value for non-existing parameter: expected empty string, got '%s'", nonExistingVal)
	}

	// Test case 3: get a parameter that is not a string
	numVal := catchPanicString(func() string { return req.GetString("num") })
	if panicValue == nil {
		t.Errorf("GetString GetString did not panic for non-string parameter")
	} else if numVal != "" {
		t.Errorf("GetString returned wrong value for non-string parameter: expected empty string, got '%s'", numVal)
	}
}

func TestGetInt(t *testing.T) {
	t.Cleanup(func() {
		panicValue = nil
	})

	// Define a test request with some parameters
	req := &Request{
		Parameters: map[string]interface{}{
			"num1": 42,
			"num2": "84",
			"str":  "not a number",
		},
	}

	// Test case 1: get an existing integer parameter
	num1Val := catchPanicInt(func() int { return req.GetInt("num1") })
	if panicValue != nil {
		t.Errorf("GetInt panic'd: %v", panicValue)
	} else if num1Val != 42 {
		t.Errorf("GetInt returned wrong value: expected 42, got %d", num1Val)
	}

	// Test case 2: get an existing parameter that is a string representation of an integer
	num2Val := catchPanicInt(func() int { return req.GetInt("num2") })
	if panicValue != nil {
		t.Errorf("GetInt panic'd: %v", panicValue)
	} else if num2Val != 84 {
		t.Errorf("GetInt returned wrong value: expected 84, got %d", num2Val)
	}

	// Test case 3: get a non-existing parameter
	nonExistingVal := catchPanicInt(func() int { return req.GetInt("non-existing") })
	if panicValue == nil {
		t.Errorf("GetInt did not panic for non-existing parameter")
	} else if nonExistingVal != 0 {
		t.Errorf("GetInt returned wrong value for non-existing parameter: expected 0, got %d", nonExistingVal)
	}

	// Test case 4: get a parameter that is not an integer
	strVal := catchPanicInt(func() int { return req.GetInt("str") })
	if panicValue == nil {
		t.Errorf("GetInt did not panic for non-integer parameter")
	} else if strVal != 0 {
		t.Errorf("GetInt returned wrong value for non-integer parameter: expected 0, got %d", strVal)
	}
}

func TestGetFloat(t *testing.T) {
	t.Cleanup(func() {
		panicValue = nil
	})

	// Define a test request with some parameters
	req := &Request{
		Parameters: map[string]interface{}{
			"float1": 3.14,
			"float2": "2.718",
			"str":    "not a number",
		},
	}

	// Test case 1: get an existing float parameter
	float1Val := catchPanicFloat(func() float64 { return req.GetFloat("float1") })
	if panicValue != nil {
		t.Errorf("GetFloat panicked: %v", panicValue)
	} else if float1Val != 3.14 {
		t.Errorf("GetFloat returned wrong value: expected 3.14, got %f", float1Val)
	}

	// Test case 2: get an existing parameter that is a string representation of a float
	float2Val := catchPanicFloat(func() float64 { return req.GetFloat("float2") })
	if panicValue != nil {
		t.Errorf("GetFloat panicked: %v", panicValue)
	} else if float2Val != 2.718 {
		t.Errorf("GetFloat returned wrong value: expected 2.718, got %f", float2Val)
	}

	// Test case 3: get a non-existing parameter
	nonExistingVal := catchPanicFloat(func() float64 { return req.GetFloat("non-existing") })
	if panicValue == nil {
		t.Errorf("GetFloat did not panic for non-existing parameter")
	} else if nonExistingVal != 0 {
		t.Errorf("GetFloat returned wrong value for non-existing parameter: expected 0, got %f", nonExistingVal)
	}

	// Test case 4: get a parameter that is not a float
	strVal := catchPanicFloat(func() float64 { return req.GetFloat("str") })
	if panicValue == nil {
		t.Errorf("GetFloat did not panic for non-float parameter")
	} else if strVal != 0 {
		t.Errorf("GetFloat returned wrong value for non-float parameter: expected 0, got %f", strVal)
	}
}

func TestGetBool(t *testing.T) {
	t.Cleanup(func() {
		panicValue = nil
	})

	// Define a test request with some parameters
	req := &Request{
		Parameters: map[string]interface{}{
			"bool1": true,
			"bool2": "false",
			"str":   "not a boolean",
		},
	}

	// Test case 1: get an existing boolean parameter
	bool1Val := catchPanicBool(func() bool { return req.GetBool("bool1") })
	if panicValue != nil {
		t.Errorf("GetBool panicked: %v", panicValue)
	} else if !bool1Val {
		t.Errorf("GetBool returned wrong value: expected true, got %t", bool1Val)
	}

	// Test case 2: get an existing parameter that is a string representation of a boolean
	bool2Val := catchPanicBool(func() bool { return req.GetBool("bool2") })
	if panicValue != nil {
		t.Errorf("GetBool panicked: %v", panicValue)
	} else if bool2Val {
		t.Errorf("GetBool returned wrong value: expected false, got %t", bool2Val)
	}

	// Test case 3: get a non-existing parameter
	nonExistingVal := catchPanicBool(func() bool { return req.GetBool("non-existing") })
	if panicValue == nil {
		t.Errorf("GetBool did not panic for non-existing parameter")
	} else if nonExistingVal != false {
		t.Errorf("GetBool returned wrong value for non-existing parameter: expected false, got %t", nonExistingVal)
	}

	// Test case 4: get a parameter that is not a boolean
	strVal := catchPanicBool(func() bool { return req.GetBool("str") })
	if panicValue == nil {
		t.Errorf("GetBool did not panic for non-boolean parameter")
	} else if strVal != false {
		t.Errorf("GetBool returned wrong value for non-boolean parameter: expected false, got %t", strVal)
	}
}

func TestGetStringArray(t *testing.T) {
	t.Cleanup(func() {
		panicValue = nil
	})

	// Define a test request with some parameters
	req := &Request{
		Parameters: map[string]interface{}{
			"array1": []int{1, 2, 3},
			"array2": []interface{}{"1", "2", "3"},
			"str":    "not an array",
		},
	}

	// Test case 1: get an existing array parameter
	array1Val := catchPanicArray(func() []int { return req.GetIntArray("array1") })
	if panicValue != nil {
		t.Errorf("GetArray panicked: %v", panicValue)
	} else if !reflect.DeepEqual(array1Val, []int{1, 2, 3}) {
		t.Errorf("GetArray returned wrong value: expected [1, 2, 3], got %v", array1Val)
	}

	// Test case 2: get an existing parameter that is an array of interfaces
	array2Val := catchPanicArray(func() []int { return req.GetIntArray("array2") })
	if panicValue != nil {
		t.Errorf("GetArray panicked: %v", panicValue)
	} else if !reflect.DeepEqual(array2Val, []int{1, 2, 3}) {
		t.Errorf("GetArray returned wrong value: expected [1, 2, 3], got %v", array2Val)
	}

	// Test case 3: get a non-existing parameter
	nonExistingVal := catchPanicArray(func() []int { return req.GetIntArray("non-existing") })
	if panicValue == nil {
		t.Errorf("GetArray did not panic for non-existing parameter")
	} else if nonExistingVal != nil {
		t.Errorf("GetArray returned wrong value for non-existing parameter: expected nil, got %v", nonExistingVal)
	}

	// Test case 4: get a parameter that is not an array
	strVal := catchPanicArray(func() []int { return req.GetIntArray("str") })
	if panicValue == nil {
		t.Errorf("GetArray did not panic for non-array parameter")
	} else if strVal != nil {
		t.Errorf("GetArray returned wrong value for non-array parameter: expected nil, got %v", strVal)
	}
}
