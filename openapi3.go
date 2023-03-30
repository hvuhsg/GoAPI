package goapi

import (
	"fmt"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hvuhsg/goapi/validators"
)

// openapi3Schema generates the OpenAPI-3 schema for the given App
func openapi3Schema(a *App) ([]byte, error) {
	paths := make(openapi3.Paths)

	// Loop through each view defined in the app
	for _, view := range a.views {
		path, ok := paths[view.path]
		if !ok {
			path = &openapi3.PathItem{}
		}

		// Loop through each HTTP method defined for the view
		for method := range view.methods {
			methodString := methodCodeToString[method]
			parameters := make(openapi3.Parameters, 0)

			// Loop through each parameter defined for the view
			for paramName, paramInfo := range view.parameters {
				schemaVal := openapi3.NewSchema()

				// Set the required field to false by default
				required := false

				// Loop through each validator defined for the parameter
				for _, validator := range paramInfo.validators {
					// Check if the validator is a VRequired validator, if so set the required field to true
					_, ok = validator.(validators.VRequired)
					if ok {
						required = true
					}

					validator.UpdateOpenAPISchema(schemaVal)
				}

				// Create a new Parameter object and add it to the Parameters array
				paramRef := openapi3.ParameterRef{Value: &openapi3.Parameter{
					Name:       paramName,
					In:         paramInfo.in,
					Required:   required,
					Schema:     openapi3.NewSchemaRef("", schemaVal),
					Deprecated: false,
				}}
				parameters = append(parameters, &paramRef)
			}

			responses := openapi3.NewResponses()

			// Set default validation error response
			respDesc := "Error when validating request against validators"
			responses["422"] = &openapi3.ResponseRef{Value: &openapi3.Response{Description: &respDesc}}

			// Create a new Operation object to hold all the information for the HTTP method
			operation := openapi3.Operation{
				Description: view.description,
				Tags:        view.tags,
				Parameters:  parameters,
				Responses:   responses,
				Deprecated:  view.depreceted,
			}

			path.SetOperation(methodString, &operation)
		}

		paths[view.path] = path
	}

	// Create the final OpenAPI-3 schema object with the Paths object and other app information
	schemaObj := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          a.title,
			Version:        a.version,
			Description:    a.description,
			License:        &a.license,
			TermsOfService: a.termOfServiceURL,
			Contact:        &a.contact,
		},
		Security: a.security,
		Tags:     a.tags,
		Paths:    paths,
	}

	// Marshal the OpenAPI-3 schema object to JSON and return it
	return schemaObj.MarshalJSON()
}

func registerDocs(a *App, mux *http.ServeMux) {
	// Docs internal view, retunrs the OpenAPI-3 schmea
	schema, schemaErr := openapi3Schema(a) // Only marshel on startup for performence
	mux.HandleFunc(a.openapiSchemaURL, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if schemaErr != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(schema)
	})

	mux.HandleFunc(a.openapiDocsURL, func(w http.ResponseWriter, _ *http.Request) {
		swaggerJsUrl := "https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui-bundle.js"
		swaggerCssUrl := "https://cdn.jsdelivr.net/npm/swagger-ui-dist@3/swagger-ui.css"
		swaggerFavIconUrl := "https://fastapi.tiangolo.com/img/favicon.png"
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
		<link type="text/css" rel="stylesheet" href="%s">
		<link rel="shortcut icon" href="%s">
		<title>%s</title>
		</head>
		<body>
		<div id="swagger-ui">
		</div>
		<script src="%s"></script>
		<script>
		const ui = SwaggerUIBundle({
			url: '%s',
			dom_id: '#swagger-ui',
			presets: [
				SwaggerUIBundle.presets.apis,
				SwaggerUIBundle.SwaggerUIStandalonePreset
			],
			layout: "BaseLayout",
			deepLinking: true,
			showExtensions: true,
			showCommonExtensions: true
		})
		</script>
		</body>
		</html>
		`, swaggerCssUrl, swaggerFavIconUrl, a.title, swaggerJsUrl, a.openapiSchemaURL)
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}
