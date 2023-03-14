package validators

import (
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsEmail struct{}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (v VIsEmail) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "email"
}
func (v VIsEmail) Validate(r *request.Request, paramName string) error {
	email := r.GetString(paramName)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("parameter %s must be a valid email address", paramName)
	}
	return nil
}
