package validators

import (
	"fmt"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsTime struct {
	// The expected format of the time string. This should be a string that specifies the expected format of the time, using the following standard Go date and time format codes: https://golang.org/pkg/time/#pkg-constants.
	// For example, to specify a time in the format "15:04:05", you can set Format to "15:04:05".
	Format string
}

func (v VIsTime) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "time"
}

func (v VIsTime) Validate(r *request.Request, paramName string) error {
	timeStr := r.GetString(paramName)

	// Parse the time string based on the specified format
	parsedTime, err := time.Parse(v.Format, timeStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid time in the format %s", paramName, v.Format)
	}

	// Check if the parsed time is valid
	if parsedTime.IsZero() {
		return fmt.Errorf("parameter %s must be a valid time in the format %s", paramName, v.Format)
	}

	return nil
}
