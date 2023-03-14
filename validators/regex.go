package validators

import (
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VRegex struct {
	Regex string
}

func (v VRegex) UpdateOpenAPISchema(schema *openapi3.Schema) {
	// Not applicable
}

func (v VRegex) Validate(r *request.Request, paramName string) error {
	str := r.GetString(paramName)

	// Compile the regular expression
	re, err := regexp.Compile(v.Regex)
	if err != nil {
		return fmt.Errorf("failed to compile regular expression: %v", err)
	}

	// Match the string against the regular expression
	if !re.MatchString(str) {
		return fmt.Errorf("parameter %s must match the regular expression %s", paramName, v.Regex)
	}

	return nil
}
