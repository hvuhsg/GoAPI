package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsArray struct{}

func (v VIsArray) UpdateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "array" }
func (VIsArray) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'array'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetArray(paramName)
	}

	return nil
}
