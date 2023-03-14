package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsBool struct{}

func (v VIsBool) UpdateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "boolean" }
func (VIsBool) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'bool'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetBool(paramName)
	}

	return nil
}
