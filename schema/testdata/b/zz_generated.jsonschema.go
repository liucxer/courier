package b

import (
	github_com_go_courier_schema_pkg_jsonschema "github.com/liucxer/courier/schema/pkg/jsonschema"
	github_com_go_courier_schema_pkg_jsonschema_extractors "github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
)

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/b.Data", new(Data))
}

func (Data) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"id": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "ID",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"id",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/b.Data",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/b.PullPolicy", new(PullPolicy))
}

func (PullPolicy) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Enum: []interface{}{
					"Always",
					"IfNotPresent",
					"Never",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-enum-labels": []string{
					"pull always",
					"if not preset",
					"never",
				},
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/b.PullPolicy",
			},
		},
	})
}
