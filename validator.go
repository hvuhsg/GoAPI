package goapi

import "fmt"

type Validator interface {
	Validate(r *Request, paramName string) error
}

type VRequired struct{}

func (VRequired) Validate(r *Request, paramName string) error {
	_, ok := r.parameters[paramName]
	if !ok {
		return fmt.Errorf("parameter %s is required", paramName)
	}

	return nil
}

type VIsString struct{}

func (VIsString) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	value := r.parameters[paramName]

	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'string'", paramName)
	}

	return nil
}

type VIsInt struct{}

func (VIsInt) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	_, ok := r.parameters[paramName].(int)
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'int'", paramName)
	}

	return nil
}

type VIsFloat struct{}

func (VIsFloat) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	_, ok := r.parameters[paramName].(float64)
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'float'", paramName)
	}

	return nil
}

type VIsBool struct{}

func (VIsBool) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	_, ok := r.parameters[paramName].(bool)
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'bool'", paramName)
	}

	return nil
}

type VIsArray struct{}

func (VIsArray) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	_, ok := r.parameters[paramName].([]interface{})
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'array'", paramName)
	}

	return nil
}

type VIsMap struct{}

func (VIsMap) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	_, ok := r.parameters[paramName].(map[string]interface{})
	if !ok {
		return fmt.Errorf("parameter %s must be of type 'map'", paramName)
	}

	return nil
}

type VStringLength struct {
	Min int
	Max int
}

func (v VStringLength) Validate(r *Request, paramName string) error {
	vr := VIsString{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	strValue, _ := r.parameters[paramName].(string)

	if len(strValue) < v.Min || len(strValue) > v.Max {
		return fmt.Errorf("parameter %s length must be between %d and %d characters", paramName, v.Min, v.Max)
	}

	return nil
}

type VRange struct {
	Min float64
	Max float64
}

func (v VRange) Validate(r *Request, paramName string) error {
	vr := VRequired{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	val := r.parameters[paramName]
	var fVal float64

	switch val := val.(type) {
	case int:
		fVal = float64(val)
	case float64:
		fVal = val
	default:
		return fmt.Errorf("parameter %s must be a numeric value", paramName)
	}

	if fVal < v.Min || fVal > v.Max {
		return fmt.Errorf("parameter %s must be between %f and %f", paramName, v.Min, v.Max)
	}

	return nil
}
