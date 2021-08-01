package extractors

import (
	"context"
	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/schema/pkg/jsonschema"
	"github.com/liucxer/courier/schema/pkg/ptr"
	"github.com/liucxer/courier/schema/pkg/validator"
)

func BindSchemaValidationByValidateBytes(s *jsonschema.Schema, typ typesutil.Type, validateBytes []byte) error {
	fieldValidator, err := validator.ValidatorMgrDefault.Compile(context.Background(), validateBytes, typ, func(rule *validator.Rule) {
		rule.DefaultValue = nil
	})
	if err != nil {
		return err
	}
	if fieldValidator != nil {
		BindSchemaValidationByValidator(s, fieldValidator)
		s.AddExtension(jsonschema.XTagValidate, string(validateBytes))
	}
	return nil
}

func BindSchemaValidationByValidator(s *jsonschema.Schema, v validator.Validator) {
	if validatorLoader, ok := v.(*validator.ValidatorLoader); ok {
		v = validatorLoader.Validator
	}

	if s == nil {
		return
	}

	switch vt := v.(type) {
	case *validator.UintValidator:
		if len(vt.Enums) > 0 {
			for v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		s.Minimum = ptr.Float64(float64(vt.Minimum))
		s.Maximum = ptr.Float64(float64(vt.Maximum))

		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum
		if vt.MultipleOf > 0 {
			s.MultipleOf = ptr.Float64(float64(vt.MultipleOf))
		}
	case *validator.IntValidator:
		if len(vt.Enums) > 0 {
			for v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		if vt.Minimum != nil {
			s.Minimum = ptr.Float64(float64(*vt.Minimum))
		}
		if vt.Maximum != nil {
			s.Maximum = ptr.Float64(float64(*vt.Maximum))
		}
		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum

		if vt.MultipleOf > 0 {
			s.MultipleOf = ptr.Float64(float64(vt.MultipleOf))
		}
	case *validator.FloatValidator:
		if len(vt.Enums) > 0 {
			for v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		if vt.Minimum != nil {
			s.Minimum = ptr.Float64(*vt.Minimum)
		}
		if vt.Maximum != nil {
			s.Maximum = ptr.Float64(*vt.Maximum)
		}
		s.ExclusiveMinimum = vt.ExclusiveMinimum
		s.ExclusiveMaximum = vt.ExclusiveMaximum

		if vt.MultipleOf > 0 {
			s.MultipleOf = ptr.Float64(vt.MultipleOf)
		}
	case *validator.StrfmtValidator:
		s.Type = []string{"string"} // force to type string for TextMarshaler
		s.Format = vt.Names()[0]
	case *validator.StringValidator:
		s.Type = []string{"string"} // force to type string for TextMarshaler

		if len(vt.Enums) > 0 {
			for v := range vt.Enums {
				s.Enum = append(s.Enum, v)
			}
			return
		}

		s.MinLength = ptr.Uint64(vt.MinLength)
		if vt.MaxLength != nil {
			s.MaxLength = ptr.Uint64(*vt.MaxLength)
		}
		if vt.Pattern != "" {
			s.Pattern = vt.Pattern
		}
	case *validator.SliceValidator:
		s.MinItems = ptr.Uint64(vt.MinItems)
		if vt.MaxItems != nil {
			s.MaxItems = ptr.Uint64(*vt.MaxItems)
		}

		if vt.ElemValidator != nil {
			if s.Items == nil {
				s.Items = &jsonschema.SchemaOrArray{}
			}

			BindSchemaValidationByValidator(s.Items.Schema, vt.ElemValidator)
		}
	case *validator.MapValidator:
		s.MinProperties = ptr.Uint64(vt.MinProperties)
		if vt.MaxProperties != nil {
			s.MaxProperties = ptr.Uint64(*vt.MaxProperties)
		}
		if vt.ElemValidator != nil {
			BindSchemaValidationByValidator(s.AdditionalProperties.Schema, vt.ElemValidator)
		}
	}
}
