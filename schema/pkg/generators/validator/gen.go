package validator

import (
	"context"
	"fmt"
	"go/ast"
	"go/types"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/liucxer/courier/gengo/pkg/gengo"
	"github.com/liucxer/courier/gengo/pkg/namer"
	gengotypes "github.com/liucxer/courier/gengo/pkg/types"
	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/schema/pkg/validator"
)

func init() {
	gengo.Register(&validatorGen{})
}

type validatorGen struct {
	imports namer.ImportTracker
}

func (g *validatorGen) Name() string {
	return "validator"
}

func (validatorGen) New() gengo.Generator {
	return &validatorGen{
		imports: namer.NewDefaultImportTracker(),
	}
}

func (g *validatorGen) Imports(context *gengo.Context) map[string]string {
	return g.imports.Imports()
}

func (g *validatorGen) Init(c *gengo.Context, writer io.Writer) error {
	return nil
}

func (g *validatorGen) GenerateType(c *gengo.Context, named *types.Named, w io.Writer) error {
	sw := gengo.NewSnippetWriter(w, namer.NameSystems{
		"raw": namer.NewRawNamer(c.Package.Pkg().Path(), g.imports),
	})

	sw.Do(`
func (v *[[ .typeName ]]) Validate() error {
	return [[ .validateFn | render ]](v)
}`, gengo.Args{
		"typeName": named.Obj().Name(),
		"validateFn": g.newValidateFn(c, typesutil.FromTType(named), &validator.ValidatorLoader{
			Validator: &validator.StructValidator{},
		}),
	})

	return nil
}

func (g *validatorGen) newValidateFn(c *gengo.Context, tt typesutil.Type, vt *validator.ValidatorLoader) func(sw *gengo.SnippetWriter) {
	d := gengo.NewDumper(c.Package.Pkg().Path(), g.imports)

	return gengo.Snippet(`
func(v *[[ .typeName ]]) error {
[[ .doValidate | render ]]  }`, gengo.Args{
		"typeName": d.TypeLit(tt),
		"doValidate": func(sw *gengo.SnippetWriter) {
			if vt.Validator == nil {
				sw.Do("return nil")
				return
			}

			defer func() {
				if e := recover(); e != nil {
					panic(fmt.Errorf("failed to gen for %s", tt))
				}
			}()

			if _, ok := typesutil.EncodingTextMarshalerTypeReplacer(tt); ok {
				sw.Do(`
vv, _ := v.MarshalText()
if vv != nil {
	return ([[ .vt ]]).Validate(string(vv));	
}
[[ if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
`, gengo.Args{
					"vt":         d.ValueLit(vt.Validator),
					"required":   !vt.Optional,
					"missingErr": d.ValueLit(&validator.MissingRequired{}),
				})
				return
			}

			switch tt.Kind() {
			case reflect.Ptr:
				sw.Do(`
vv := *v
if vv == nil {
	[[ if .shouldSetDefault ]] var vvv [[ .typeName ]]; vv = &vvv; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return [[ .validateFn | render ]](vv);	
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"required":         !vt.Optional,
					"validateFn":       g.newValidateFn(c, tt.Elem(), vt),
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
					"typeName":         d.TypeLit(tt.Elem()),
				})
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(int64(vv));	
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(uint64(vv));	
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Float64, reflect.Float32:
				sw.Do(`
vv := *v
if vv == 0 {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(float64(vv));
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     string(vt.DefaultValue),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.String:
				sw.Do(`
vv := *v
if vv == "" {
	[[ if .shouldSetDefault ]] vv = [[ .defaultValue ]]; *v = vv; [[ else if .required ]] return [[ .missingErr ]] [[ else ]] return nil [[ end ]]
}
return ([[ .vt ]]).Validate(string(vv));
`, gengo.Args{
					"shouldSetDefault": len(vt.DefaultValue) != 0,
					"defaultValue":     strconv.Quote(string(vt.DefaultValue)),
					"vt":               d.ValueLit(vt.Validator),
					"required":         !vt.Optional,
					"missingErr":       d.ValueLit(&validator.MissingRequired{}),
				})
			case reflect.Map:
				mapValidator := vt.Validator.(*validator.MapValidator)

				sw.Do(`
m := *v;
`)

				if mapValidator.MinProperties > 0 {
					sw.Do(`
if n := len(m); n < [[ .minProperties ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"minProperties": mapValidator.MinProperties,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Minimum: mapValidator.MinProperties,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				if mapValidator.MaxProperties != nil {
					sw.Do(`
if n := len(m); n > [[ .maxItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"maxItems": *mapValidator.MaxProperties,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Maximum: *mapValidator.MaxProperties,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				args := gengo.Args{
					"FnNewErrSet":       gengotypes.Ref("github.com/liucxer/courier/schema/pkg/validator", "NewErrorSet"),
					"validateItemFn":    g.newValidateFn(c, tt.Elem(), mapValidator.ElemValidator.(*validator.ValidatorLoader)),
					"shouldValidateKey": mapValidator.KeyValidator != nil,
				}

				if mapValidator.KeyValidator != nil {
					args["validateKeyFn"] = g.newValidateFn(c, tt.Key(), mapValidator.KeyValidator.(*validator.ValidatorLoader))
				}

				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()

[[ if .shouldValidateKey ]] validateKey := [[ .validateKeyFn | render ]] [[ end ]]
validateElem := [[ .validateItemFn | render ]]

for k := range m {
	[[ if .shouldValidateKey ]] if e := validateKey(&k); e != nil {
			errSet.AddErr(e, k)
	} [[ end ]]
	value := m[k]
	if e := validateElem(&value); e != nil {
		errSet.AddErr(e, k)
	}
}

if errSet.Len() > 0 {
	return errSet
}
return nil
`, args)

			case reflect.Array, reflect.Slice:
				sliceValidator := vt.Validator.(*validator.SliceValidator)

				sw.Do(`
list := *v;
n := len(list);
`)

				if sliceValidator.MinItems > 0 {
					sw.Do(`
if n < [[ .minItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"minItems": sliceValidator.MinItems,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Minimum: sliceValidator.MinItems,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				if sliceValidator.MaxItems != nil {
					sw.Do(`
if n > [[ .maxItems ]] {
	return [[ .err ]]
}
`, gengo.Args{
						"maxItems": *sliceValidator.MaxItems,
						"err": d.ValueLit(&validator.OutOfRangeError{
							Current: "CURRENT",
							Maximum: *sliceValidator.MaxItems,
							Target:  validator.TargetSliceLength,
						}, gengo.OnInterface(func(v interface{}) string {
							if v == "CURRENT" {
								return "n"
							}
							return fmt.Sprintf("%v", v)
						})),
					})
				}

				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()

validateItem := [[ .validateItemFn | render ]]

for i := 0; i < n; i++ {
	errSet.AddErr(validateItem(&list[i]), i)
}

if errSet.Len() > 0 {
	return errSet
}
return nil
`, gengo.Args{
					"FnNewErrSet":    gengotypes.Ref("github.com/liucxer/courier/schema/pkg/validator", "NewErrorSet"),
					"validateItemFn": g.newValidateFn(c, tt.Elem(), sliceValidator.ElemValidator.(*validator.ValidatorLoader)),
				})

			case reflect.Struct:
				sw.Do(`
errSet := [[ .FnNewErrSet | raw ]]()
[[ .validateFields | render ]]
if errSet.Len() > 0 {
	return errSet
}
return nil
`, gengo.Args{
					"FnNewErrSet": gengotypes.Ref("github.com/liucxer/courier/schema/pkg/validator", "NewErrorSet"),
					"validateFields": func(sw *gengo.SnippetWriter) {
						walkStruct(tt, func(sf *StructField) {
							if sf.Validator == nil {
								return
							}

							sw.Do(`
errSet.AddErr([[ .doValidateFn | render ]](&v.[[ .fieldName ]]), [[ if .hasLocation ]] [[ .location ]], [[ end ]] [[ .displayName ]])
`, gengo.Args{
								"fieldName":    sf.Name,
								"displayName":  strconv.Quote(sf.DisplayName),
								"doValidateFn": g.newValidateFn(c, sf.Type, sf.Validator),
								"hasLocation":  sf.Location != "",
								"location":     fmt.Sprintf("%s(%s)", d.ReflectTypeLit(reflect.TypeOf(sf.Location)), d.ValueLit(sf.Location)),
							})
						}, []string{})
					},
				})
			}
		},
	})
}

type StructField struct {
	Name        string
	DisplayName string
	Location    validator.Location
	Type        typesutil.Type
	Validator   *validator.ValidatorLoader
}

func walkStruct(s typesutil.Type, each func(sf *StructField), parents []string) {
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fieldName := field.Name()
		fieldTag := field.Tag()

		displayNameTag, hasDisplayNameTag := fieldTag.Lookup("json")
		if !hasDisplayNameTag {
			displayNameTag, hasDisplayNameTag = fieldTag.Lookup("name")
		}

		locationTag, _ := fieldTag.Lookup("in")

		validateTag, _ := fieldTag.Lookup("validate")

		defaultValue, hasDefaultValue := fieldTag.Lookup("default")

		displayName := strings.Split(displayNameTag, ",")[0]
		omitempty := strings.Contains(displayNameTag, "omitempty")

		if !ast.IsExported(fieldName) || displayName == "-" {
			continue
		}

		fieldType := field.Type()

		for {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			} else {
				break
			}
		}

		if field.Anonymous() {
			switch fieldType.Kind() {
			case reflect.Struct:
				if !hasDisplayNameTag {
					walkStruct(fieldType, each, parents)
					continue
				}
			case reflect.Interface:
				continue
			}
		}

		sf := &StructField{}

		sf.Name = fieldName
		sf.DisplayName = displayName
		sf.Type = field.Type()
		sf.Location = validator.Location(locationTag)

		if validateTag != "-" {
			v := validator.ValidatorMgrDefault.MustCompile(context.Background(), []byte(validateTag), fieldType, func(rule *validator.Rule) {
				rule.Optional = omitempty

				if hasDefaultValue {
					rule.DefaultValue = []byte(defaultValue)
				}
			})

			vl := v.(*validator.ValidatorLoader)
			if vl != nil && vl.Validator != nil {
				sf.Validator = vl
			}
		}

		each(sf)
	}
}
