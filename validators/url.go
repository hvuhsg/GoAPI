package validators

import (
	"fmt"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/request"
)

type VIsURL struct{}

func (v VIsURL) UpdateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "url"
}
func (v VIsURL) Validate(r *request.Request, paramName string) error {
	urlStr := r.GetString(paramName)

	// Parse the URL and check for any errors
	_, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid URL: %v", paramName, err)
	}

	return nil
}
