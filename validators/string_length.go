package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VStringLength struct {
	Min int
	Max int
}

func (v VStringLength) UpdateOpenAPISchema(schema *openapi3.Schema) {
	VIsString{}.UpdateOpenAPISchema(schema)
	schema.MinLength = uint64(v.Min)
	maxL := uint64(v.Max)
	schema.MaxLength = &maxL
}
func (v VStringLength) Validate(r *request.Request, paramName string) error {
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
