package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsString struct{}

func (v VIsString) UpdateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "string" }
func (VIsString) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'string'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetString(paramName)
	}

	return nil
}
