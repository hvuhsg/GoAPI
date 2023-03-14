package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VRange struct {
	Min float64
	Max float64
}

func (v VRange) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Min = &v.Min
	schema.Max = &v.Max
}
func (v VRange) Validate(r *request.Request, paramName string) error {
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
