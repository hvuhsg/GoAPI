package goapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Request struct {
	HTTPRequest *http.Request
	parameters  map[string]any
}

func NewRequest(req *http.Request) (*Request, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	params := make(map[string]interface{})
	for k, v := range req.Form {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	for k, v := range req.URL.Query() {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	for k, v := range req.PostForm {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	decoder := json.NewDecoder(req.Body)
	var bodyParams interface{}
	err = decoder.Decode(&bodyParams)
	if err != nil && err != io.EOF {
		return nil, err
	}

	for k, v := range bodyParams.(map[string]interface{}) {
		params[k] = v
	}

	return &Request{
		HTTPRequest: req,
		parameters:  params,
	}, nil
}

func (r *Request) GetString(name string) (string, error) {
	val, ok := r.parameters[name]
	if !ok {
		return "", fmt.Errorf("parameter '%s' not found", name)
	}

	sVal, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("parameter '%s' is not a string", name)
	}

	return sVal, nil
}

func (r *Request) GetInt(name string) (int, error) {
	val, ok := r.parameters[name]
	if !ok {
		return 0, fmt.Errorf("parameter '%s' not found", name)
	}

	switch v := val.(type) {
	case int:
		return v, nil
	case string:
		iVal, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("parameter '%s' is not an integer: %s", name, err)
		}
		return iVal, nil
	default:
		return 0, fmt.Errorf("parameter '%s' is not an integer", name)
	}
}

func (r *Request) GetFloat(name string) (float64, error) {
	val, ok := r.parameters[name]
	if !ok {
		return 0, fmt.Errorf("parameter '%s' not found", name)
	}

	switch v := val.(type) {
	case float64:
		return v, nil
	case string:
		fVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, fmt.Errorf("parameter '%s' is not a float: %s", name, err)
		}
		return fVal, nil
	default:
		return 0, fmt.Errorf("parameter '%s' is not a float", name)
	}
}

func (r *Request) GetBool(name string) (bool, error) {
	val, ok := r.parameters[name]
	if !ok {
		return false, fmt.Errorf("parameter '%s' not found", name)
	}

	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("parameter '%s' is not a boolean: %s", name, err)
		}
		return bVal, nil
	default:
		return false, fmt.Errorf("parameter '%s' is not a boolean", name)
	}
}

func (r *Request) GetArray(name string) ([]interface{}, error) {
	value, ok := r.parameters[name].([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter %s not found or not an array", name)
	}
	return value, nil
}

func (r *Request) GetMap(name string) (map[string]interface{}, error) {
	value, ok := r.parameters[name].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter %s not found or not a map", name)
	}
	return value, nil
}

func (r *Request) GetIntArray(name string) ([]int, error) {
	val, ok := r.parameters[name]
	if !ok {
		return nil, fmt.Errorf("parameter '%s' not found", name)
	}

	arr, ok := val.([]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter '%s' is not an array", name)
	}

	intArr := make([]int, len(arr))
	for i, v := range arr {
		intVal, ok := v.(int)
		if !ok {
			return nil, fmt.Errorf("parameter '%s' contains a non-integer value", name)
		}
		intArr[i] = intVal
	}

	return intArr, nil
}

func (r *Request) GetStringBoolMap(name string) (map[string]bool, error) {
	val, ok := r.parameters[name]
	if !ok {
		return nil, fmt.Errorf("parameter '%s' not found", name)
	}

	m, ok := val.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter '%s' is not a map", name)
	}

	boolMap := make(map[string]bool)
	for k, v := range m {
		bVal, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("parameter '%s' contains a non-boolean value", name)
		}
		boolMap[k] = bVal
	}

	return boolMap, nil
}
