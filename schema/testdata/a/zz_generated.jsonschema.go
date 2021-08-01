package a

import (
	github_com_go_courier_schema_pkg_jsonschema "github.com/liucxer/courier/schema/pkg/jsonschema"
	github_com_go_courier_schema_pkg_jsonschema_extractors "github.com/liucxer/courier/schema/pkg/jsonschema/extractors"
)

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.ArrayString", new(ArrayString))
}

func (ArrayString) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"array",
			},
			Items: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrArray{
				Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
				}),
			}),
			Description: "ArrayString",
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				MaxItems: func(v uint64) *uint64 { return &v }(2),
				MinItems: func(v uint64) *uint64 { return &v }(2),
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.ArrayString",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Bool", new(Bool))
}

func (Bool) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"boolean",
			},
			Description: "Bool",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Bool",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Composed", new(Composed))
}

func (Composed) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/liucxer/courier/schema/testdata/a.Part"),
					},
				}),
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"object",
						},
					},
				}),
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Composed",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Double", new(Double))
}

func (Double) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"number",
			},
			Format:      "double",
			Description: "Double",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Double",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Float", new(Float))
}

func (Float) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"number",
			},
			Format:      "float",
			Description: "Float",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Float",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Int", new(Int))
}

func (Int) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"integer",
			},
			Format:      "int32",
			Description: "Int",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Int",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.InvalidComposed", new(InvalidComposed))
}

func (InvalidComposed) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/liucxer/courier/schema/testdata/a.Part"),
					},
				}),
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/liucxer/courier/schema/testdata/a.PartConflict"),
					},
				}),
				&(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"object",
						},
					},
				}),
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.InvalidComposed",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Map", new(Map))
}

func (Map) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			AdditionalProperties: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrBool{
				Allows: true,
				Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/liucxer/courier/schema/testdata/a.String"),
					},
				}),
			}),
			PropertyNames: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
				SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
					Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
						"string",
					},
				},
			}),
			Description: "Map",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Map",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.NamedComposed", new(NamedComposed))
}

func (NamedComposed) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"part": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/liucxer/courier/schema/testdata/a.Part"),
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "Part",
									},
								},
							}),
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"part",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.NamedComposed",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Node", new(Node))
}

func (Node) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"children": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"array",
						},
						Items: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrArray{
							Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Nullable: true,
								},
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/liucxer/courier/schema/testdata/a.Node"),
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-star-level": 1,
									},
								},
							}),
						}),
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Children",
						},
					},
				}),
				"type": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Type",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"type",
					"children",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Node",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Part", new(Part))
}

func (Part) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"Name": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MinLength: func(v uint64) *uint64 { return &v }(2),
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Name",
							"x-tag-validate":  "@string[2,]",
						},
					},
				}),
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Part",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.PartConflict", new(PartConflict))
}

func (PartConflict) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"object",
			},
			Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
				"name": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MinLength: func(v uint64) *uint64 { return &v }(0),
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Name",
							"x-tag-validate":  "@string[0,]",
						},
					},
				}),
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"name",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.PartConflict",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Protocol", new(Protocol))
}

func (Protocol) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Enum: []interface{}{
					"HTTP",
					"HTTPS",
					"TCP",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-enum-labels": []string{
					"http",
					"https",
					"tcp",
				},
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.SliceNamed", new(SliceNamed))
}

func (SliceNamed) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"array",
			},
			Items: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrArray{
				Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
						Refer: ref("github.com/liucxer/courier/schema/testdata/a.String"),
					},
				}),
			}),
			Description: "SliceNamed",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.SliceNamed",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.SliceString", new(SliceString))
}

func (SliceString) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"array",
			},
			Items: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrArray{
				Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
					},
				}),
			}),
			Description: "SliceString",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.SliceString",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.String", new(String))
}

func (String) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"string",
			},
			Description: "String",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.String",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Struct", new(Struct))
}

func (Struct) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
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
						Description: "id",
						Nullable:    true,
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MinLength: func(v uint64) *uint64 { return &v }(0),
							Pattern:   "\\d+",
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "ID",
							"x-go-star-level": 2,
							"x-tag-validate":  "@string/\\d+/",
						},
					},
				}),
				"int": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"integer",
						},
						Format: "int32",
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							Maximum:          func(v float64) *float64 { return &v }(1024),
							ExclusiveMaximum: true,
							Minimum:          func(v float64) *float64 { return &v }(-2147483648),
							ExclusiveMinimum: true,
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Int",
							"x-tag-validate":  "@int(,1024)",
						},
					},
				}),
				"map": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"object",
						},
						AdditionalProperties: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrBool{
							Allows: true,
							Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
										"object",
									},
									AdditionalProperties: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrBool{
										Allows: true,
										Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
											SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
												Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
													"object",
												},
												Properties: map[string]*github_com_go_courier_schema_pkg_jsonschema.Schema{
													"id": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
														SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
															Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
																"integer",
															},
															Format: "int32",
															SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
																Maximum: func(v float64) *float64 { return &v }(10),
																Minimum: func(v float64) *float64 { return &v }(0),
															},
														},
														VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
															Extensions: map[string]interface{}{
																"x-go-field-name": "ID",
																"x-tag-validate":  "@int[0,10]",
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
										}),
									}),
									PropertyNames: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
										SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
											Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
												"string",
											},
										},
									}),
									SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
										MinProperties: func(v uint64) *uint64 { return &v }(0),
									},
								},
							}),
						}),
						PropertyNames: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
							SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
								Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
									"string",
								},
							},
						}),
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MaxProperties: func(v uint64) *uint64 { return &v }(3),
							MinProperties: func(v uint64) *uint64 { return &v }(0),
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Map",
							"x-tag-validate":  "@map<,@map<,@struct>>[0,3]",
						},
					},
				}),
				"name": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"string",
						},
						Description: "name",
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MinLength: func(v uint64) *uint64 { return &v }(2),
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Name",
							"x-go-star-level": 1,
							"x-tag-validate":  "@string[2,]",
						},
					},
				}),
				"protocol": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
										"string",
									},
									SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
										Enum: []interface{}{
											0,
										},
									},
								},
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/liucxer/courier/schema/testdata/a.Protocol"),
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-tag-validate": "@string{HTTP}",
									},
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Description: "pull policy",
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "Protocol",
									},
								},
							}),
						},
					},
				}),
				"pullPolicy": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
										"string",
									},
									SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
										Enum: []interface{}{
											0,
										},
									},
								},
								Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
									Refer: ref("github.com/liucxer/courier/schema/testdata/b.PullPolicy"),
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-tag-validate": "@string{Always}",
									},
								},
							}),
							&(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Description: "pull policy",
								},
								VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
									Extensions: map[string]interface{}{
										"x-go-field-name": "PullPolicy",
									},
								},
							}),
						},
					},
				}),
				"slice": &(github_com_go_courier_schema_pkg_jsonschema.Schema{
					SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
						Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
							"array",
						},
						Items: &(github_com_go_courier_schema_pkg_jsonschema.SchemaOrArray{
							Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
								SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
									Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
										"integer",
									},
									Format: "int32",
									SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
										Maximum: func(v float64) *float64 { return &v }(10),
										Minimum: func(v float64) *float64 { return &v }(0),
									},
								},
							}),
						}),
						SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
							MaxItems: func(v uint64) *uint64 { return &v }(30),
							MinItems: func(v uint64) *uint64 { return &v }(1),
						},
					},
					VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
						Extensions: map[string]interface{}{
							"x-go-field-name": "Slice",
							"x-tag-validate":  "@slice<@int[0,10]>[1,30]",
						},
					},
				}),
			},
			Description: "Struct",
			SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
				Required: []string{
					"name",
					"pullPolicy",
					"protocol",
					"slice",
				},
			},
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Struct",
			},
		},
	})
}

func init() {
	github_com_go_courier_schema_pkg_jsonschema_extractors.Register("github.com/liucxer/courier/schema/testdata/a.Uint", new(Uint))
}

func (Uint) OpenAPISchema(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_jsonschema.Schema {
	return &(github_com_go_courier_schema_pkg_jsonschema.Schema{
		SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
			Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
				"integer",
			},
			Format:      "uint32",
			Description: "Uint",
		},
		VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
			Extensions: map[string]interface{}{
				"x-go-vendor-type": "github.com/liucxer/courier/schema/testdata/a.Uint",
			},
		},
	})
}
