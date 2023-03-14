package validators

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VPassword struct {
	MinLength      int
	MaxLength      int
	RequireSymbols bool
	RequireNumbers bool
	RequireUpper   bool
}

func (v VPassword) UpdateOpenAPISchema(schema *openapi3.Schema) {
	// Not applicable
}

func (v VPassword) Validate(r *request.Request, paramName string) error {
	password := r.GetString(paramName)

	// Check length requirements
	if len(password) < v.MinLength || len(password) > v.MaxLength {
		return fmt.Errorf("parameter %s must have length between %d and %d", paramName, v.MinLength, v.MaxLength)
	}

	// Check symbol requirements
	if v.RequireSymbols && !regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\{\}\[\]:;<>\?,\.\|\\\/~]`).MatchString(password) {
		return errors.New("password must contain at least one symbol")
	}

	// Check number requirements
	if v.RequireNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	// Check uppercase letter requirements
	if v.RequireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	return nil
}
