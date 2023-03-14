package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsInt struct{}

func (v VIsInt) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "integer"
	schema.Format = "int64"
}
func (VIsInt) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'int'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetInt(paramName)
	}

	return nil
}
