package goapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Request struct {
	HTTPRequest *http.Request
	parameters  map[string]any
}

func NewRequest(req *http.Request) *Request {
	params := make(map[string]interface{})

	// Parse form params
	err := req.ParseForm()
	if err == nil {
		for k, v := range req.Form {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}
	}

	// Parse query params
	for k, v := range req.URL.Query() {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	// Parse post-form params
	for k, v := range req.PostForm {
		if len(v) == 1 {
			params[k] = v[0]
		} else {
			params[k] = v
		}
	}

	// Parse body params
	decoder := json.NewDecoder(req.Body)
	var bodyParams interface{}
	err = decoder.Decode(&bodyParams)
	if err == nil {
		mapBodyParams, ok := bodyParams.(map[string]interface{})
		if ok {
			for k, v := range mapBodyParams {
				params[k] = v
			}
		}
	}

	return &Request{
		HTTPRequest: req,
		parameters:  params,
	}
}

func (r *Request) GetString(name string) string {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	sVal, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("parameter '%s' is not a string", name))
	}

	return sVal
}

func (r *Request) GetInt(name string) int {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case int:
		return v
	case string:
		iVal, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("parameter '%s' is not an integer: %s", name, err))
		}
		return iVal
	default:
		panic(fmt.Sprintf("parameter '%s' is not an integer", name))
	}
}

func (r *Request) GetFloat(name string) float64 {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		fVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(fmt.Sprintf("parameter '%s' is not a float: %s", name, err))
		}
		return fVal
	default:
		panic(fmt.Sprintf("parameter '%s' is not a float", name))
	}
}

func (r *Request) GetBool(name string) bool {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case bool:
		return v
	case string:
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("parameter '%s' is not a boolean: %s", name, err))
		}
		return bVal
	default:
		panic(fmt.Sprintf("parameter '%s' is not a boolean", name))
	}
}

func (r *Request) GetArray(name string) []interface{} {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	fmt.Printf("%t\n", val)
	switch v := val.(type) {
	case []interface{}:
		return v
	case []string:
		anyArr := make([]any, len(v))
		anyArr = append(anyArr, v)
		return anyArr
	default:
		panic(fmt.Sprintf("parameter '%s' is not an array", name))
	}
}

func (r *Request) GetStringArray(name string) []string {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case []string:
		return v
	case string:
		return []string{v}
	default:
		panic(fmt.Sprintf("parameter '%s' is not an array of strings", name))
	}
}

func (r *Request) GetIntArray(name string) []int {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case []int:
		return v
	case []string:
		var intValues []int
		for _, s := range v {
			iVal, err := strconv.Atoi(s)
			if err != nil {
				panic(fmt.Sprintf("parameter '%s' contains non-integer value: %s", name, s))
			}
			intValues = append(intValues, iVal)
		}
		return intValues
	case []interface{}:
		var intValues []int
		for _, s := range v {
			sVal, ok := s.(string)
			if !ok {
				panic(fmt.Sprintf("parameter '%s' contains non-integer value: %s", name, s))
			}

			iVal, err := strconv.Atoi(sVal)
			if err != nil {
				panic(fmt.Sprintf("parameter '%s' contains non-integer value: %s", name, s))
			}

			intValues = append(intValues, iVal)
		}
		return intValues
	default:
		panic(fmt.Sprintf("parameter '%s' is not an array of integers", name))
	}
}

func (r *Request) GetFloatArray(name string) []float64 {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case []float64:
		return v
	case []string:
		var floatValues []float64
		for _, s := range v {
			fVal, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(fmt.Sprintf("parameter '%s' contains non-float value: %s", name, s))
			}
			floatValues = append(floatValues, fVal)
		}
		return floatValues
	default:
		panic(fmt.Sprintf("parameter '%s' is not an array of floats", name))
	}
}

func (r *Request) GetBoolArray(name string) []bool {
	val, ok := r.parameters[name]
	if !ok {
		panic(fmt.Sprintf("parameter '%s' not found", name))
	}

	switch v := val.(type) {
	case []bool:
		return v
	case []string:
		var boolValues []bool
		for _, s := range v {
			bVal, err := strconv.ParseBool(s)
			if err != nil {
				panic(fmt.Sprintf("parameter '%s' contains non-boolean value: %s", name, s))
			}
			boolValues = append(boolValues, bVal)
		}
		return boolValues
	default:
		panic(fmt.Sprintf("parameter '%s' is not an array of booleans", name))
	}
}

func (r *Request) GetMap(name string) (map[string]interface{}, error) {
	// FIXME
	value, ok := r.parameters[name].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("parameter %s not found or not a map", name)
	}
	return value, nil
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
