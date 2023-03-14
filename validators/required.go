package validators

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VRequired struct{}

func (v VRequired) UpdateOpenAPISchema(schema *openapi3.Schema) {}
func (VRequired) Validate(r *request.Request, paramName string) error {
	_, ok := r.Parameters[paramName]
	if !ok {
		return fmt.Errorf("parameter %s is required", paramName)
	}

	return nil
}
