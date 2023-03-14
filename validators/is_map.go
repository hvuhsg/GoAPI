package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsMap struct{}

func (v VIsMap) UpdateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "object" }
func (VIsMap) Validate(r *request.Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'map'", paramName)
		}
	}()

	_, ok := r.Parameters[paramName]

	if ok {
		r.GetMap(paramName)
	}

	return nil
}
