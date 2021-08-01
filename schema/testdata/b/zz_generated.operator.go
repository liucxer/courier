package b

import (
	github_com_go_courier_courier "github.com/liucxer/courier/courier"
	github_com_go_courier_schema_pkg_jsonschema "github.com/liucxer/courier/schema/pkg/jsonschema"
	github_com_go_courier_schema_pkg_openapi "github.com/liucxer/courier/schema/pkg/openapi"
)

func (GetByID) New() github_com_go_courier_courier.Operator {
	return &GetByID{}
}

func (GetByID) OpenAPIOperation(ref func(t string) github_com_go_courier_schema_pkg_jsonschema.Refer) *github_com_go_courier_schema_pkg_openapi.Operation {
	return &(github_com_go_courier_schema_pkg_openapi.Operation{
		OperationObject: github_com_go_courier_schema_pkg_openapi.OperationObject{
			OperationId: "GetByID",
			Parameters: []*github_com_go_courier_schema_pkg_openapi.Parameter{
				&(github_com_go_courier_schema_pkg_openapi.Parameter{
					ParameterObject: github_com_go_courier_schema_pkg_openapi.ParameterObject{
						Name: "id",
						In:   "path",
						ParameterCommonObject: github_com_go_courier_schema_pkg_openapi.ParameterCommonObject{
							Required: true,
							WithContentOrSchema: github_com_go_courier_schema_pkg_openapi.WithContentOrSchema{
								Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
									SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
										Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
											"string",
										},
										Description: "ID",
									},
									VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
										Extensions: map[string]interface{}{
											"x-go-field-name": "ID",
										},
									},
								}),
							},
						},
					},
				}),
				&(github_com_go_courier_schema_pkg_openapi.Parameter{
					ParameterObject: github_com_go_courier_schema_pkg_openapi.ParameterObject{
						Name: "pullPolicy",
						In:   "query",
						ParameterCommonObject: github_com_go_courier_schema_pkg_openapi.ParameterCommonObject{
							WithContentOrSchema: github_com_go_courier_schema_pkg_openapi.WithContentOrSchema{
								Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
									SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
										AllOf: []*github_com_go_courier_schema_pkg_jsonschema.Schema{
											&(github_com_go_courier_schema_pkg_jsonschema.Schema{
												Reference: github_com_go_courier_schema_pkg_jsonschema.Reference{
													Refer: ref("github.com/liucxer/courier/schema/testdata/b.PullPolicy"),
												},
											}),
											&(github_com_go_courier_schema_pkg_jsonschema.Schema{
												SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
													Description: "PullPolicy",
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
							},
						},
					},
				}),
				&(github_com_go_courier_schema_pkg_openapi.Parameter{
					ParameterObject: github_com_go_courier_schema_pkg_openapi.ParameterObject{
						Name: "name",
						In:   "query",
						ParameterCommonObject: github_com_go_courier_schema_pkg_openapi.ParameterCommonObject{
							WithContentOrSchema: github_com_go_courier_schema_pkg_openapi.WithContentOrSchema{
								Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
									SchemaBasic: github_com_go_courier_schema_pkg_jsonschema.SchemaBasic{
										Type: github_com_go_courier_schema_pkg_jsonschema.StringOrArray{
											"string",
										},
										Description: "Name",
										SchemaValidation: github_com_go_courier_schema_pkg_jsonschema.SchemaValidation{
											MinLength: func(v uint64) *uint64 { return &v }(4),
										},
									},
									VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
										Extensions: map[string]interface{}{
											"x-go-field-name": "Name",
											"x-tag-validate":  "@string[4,]",
										},
									},
								}),
							},
						},
					},
				}),
				&(github_com_go_courier_schema_pkg_openapi.Parameter{
					ParameterObject: github_com_go_courier_schema_pkg_openapi.ParameterObject{
						Name: "label",
						In:   "query",
						ParameterCommonObject: github_com_go_courier_schema_pkg_openapi.ParameterCommonObject{
							WithContentOrSchema: github_com_go_courier_schema_pkg_openapi.WithContentOrSchema{
								Schema: &(github_com_go_courier_schema_pkg_jsonschema.Schema{
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
										Description: "Label",
									},
									VendorExtensible: github_com_go_courier_schema_pkg_jsonschema.VendorExtensible{
										Extensions: map[string]interface{}{
											"x-go-field-name": "Label",
										},
									},
								}),
							},
						},
					},
				}),
			},
		},
	})
}
