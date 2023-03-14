package validators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsUUID struct{}

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func (v VIsUUID) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "uuid"
}
func (v VIsUUID) Validate(r *request.Request, paramName string) error {
	uuidStr := r.GetString(paramName)

	// Check if the UUID string matches the expected format
	if !uuidRegex.MatchString(uuidStr) {
		return fmt.Errorf("parameter %s must be a valid UUID", paramName)
	}

	// Convert the UUID string to lowercase and check for any non-hexadecimal characters
	uuidStr = strings.ReplaceAll(strings.ToLower(uuidStr), "-", "")
	if _, err := fmt.Sscanf(uuidStr, "%x", &struct{}{}); err != nil {
		return fmt.Errorf("parameter %s must be a valid UUID", paramName)
	}

	return nil
}
