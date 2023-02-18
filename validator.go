package goapi

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type Validator interface {
	Validate(r *Request, paramName string) error
	updateOpenAPISchema(schema *openapi3.Schema)
}

type VRequired struct{}

func (v VRequired) updateOpenAPISchema(schema *openapi3.Schema) {}
func (VRequired) Validate(r *Request, paramName string) error {
	_, ok := r.parameters[paramName]
	if !ok {
		return fmt.Errorf("parameter %s is required", paramName)
	}

	return nil
}

type VIsString struct{}

func (v VIsString) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "string" }
func (VIsString) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'string'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetString(paramName)
	}

	return nil
}

type VIsInt struct{}

func (v VIsInt) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "integer"
	schema.Format = "int64"
}
func (VIsInt) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'int'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetInt(paramName)
	}

	return nil
}

type VIsFloat struct{}

func (v VIsFloat) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "number"
	schema.Format = "float"
}
func (VIsFloat) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'float'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetFloat(paramName)
	}

	return nil
}

type VIsBool struct{}

func (v VIsBool) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "boolean" }
func (VIsBool) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'bool'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetBool(paramName)
	}

	return nil
}

type VIsArray struct{}

func (v VIsArray) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "array" }
func (VIsArray) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'array'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetArray(paramName)
	}

	return nil
}

type VIsMap struct{}

func (v VIsMap) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "object" }
func (VIsMap) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'map'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetMap(paramName)
	}

	return nil
}

type VStringLength struct {
	Min int
	Max int
}

func (v VStringLength) updateOpenAPISchema(schema *openapi3.Schema) {
	VIsString{}.updateOpenAPISchema(schema)
	schema.MinLength = uint64(v.Min)
	maxL := uint64(v.Max)
	schema.MaxLength = &maxL
}
func (v VStringLength) Validate(r *Request, paramName string) error {
	vr := VIsString{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	strValue := r.GetString(paramName)

	if len(strValue) < v.Min || len(strValue) > v.Max {
		return fmt.Errorf("parameter %s length must be between %d and %d characters", paramName, v.Min, v.Max)
	}

	return nil
}

type VRange struct {
	Min float64
	Max float64
}

func (v VRange) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Min = &v.Min
	schema.Max = &v.Max
}
func (v VRange) Validate(r *Request, paramName string) error {
	vr := VIsFloat{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	fVal := r.GetFloat(paramName)

	if fVal < v.Min || fVal > v.Max {
		return fmt.Errorf("parameter %s must be between %f and %f", paramName, v.Min, v.Max)
	}

	return nil
}
