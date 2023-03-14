package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsFloat struct{}

func (v VIsFloat) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "number"
	schema.Format = "float"
}
func (VIsFloat) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'float'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetFloat(paramName)
	}

	return nil
}
