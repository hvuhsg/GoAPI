package validators

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type Validator interface {
	Validate(r *request.Request, paramName string) error
	UpdateOpenAPISchema(schema *openapi3.Schema)
}
