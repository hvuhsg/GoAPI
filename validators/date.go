package validators

import (
	"fmt"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsDate struct {
	// The expected date format. This should be a string that specifies the expected format of the date, using the following standard Go date and time format codes: https://golang.org/pkg/time/#pkg-constants.
	// For example, to specify a date in the format "YYYY-MM-DD", you can set Format to "2006-01-02".
	Format string
}

func (v VIsDate) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "date"
}
func (v VIsDate) Validate(r *request.Request, paramName string) error {
	dateStr := r.GetString(paramName)

	// Parse the date string based on the specified format
	parsedDate, err := time.Parse(v.Format, dateStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid date in the format %s", paramName, v.Format)
	}

	// Check if the parsed date is valid
	if parsedDate.IsZero() {
		return fmt.Errorf("parameter %s must be a valid date in the format %s", paramName, v.Format)
	}

	return nil
}
