package openapi

import "github.com/liucxer/courier/schema/pkg/jsonschema"

func NewOpenAPI() *OpenAPI {
	openAPI := &OpenAPI{}
	openAPI.OpenAPI = "3.0.3"
	openAPI.Paths.Paths = map[string]*PathItem{}
	return openAPI
}

type OpenAPI struct {
	OpenAPIObject
	SpecExtensions
}

func (i OpenAPI) MarshalJSON() ([]byte, error) {
	return jsonschema.FlattenMarshalJSON(i.OpenAPIObject, i.SpecExtensions)
}

func (i *OpenAPI) UnmarshalJSON(data []byte) error {
	return jsonschema.FlattenUnmarshalJSON(data, &i.OpenAPIObject, &i.SpecExtensions)
}

type OpenAPIObject struct {
	OpenAPI string `json:"openapi"`
	Info    `json:"info"`
	Paths   `json:"paths"`
	WithServers
	WithSecurityRequirement
	WithTags
	Components `json:"components"`
}
