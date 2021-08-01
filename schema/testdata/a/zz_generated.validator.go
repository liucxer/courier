package a

import (
	github_com_go_courier_schema_pkg_validator "github.com/liucxer/courier/schema/pkg/validator"
	github_com_go_courier_schema_testdata_b "github.com/liucxer/courier/schema/testdata/b"
)

func (v *Struct) Validate() error {
	return func(v *Struct) error {
		errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
		errSet.AddErr(func(v *int) error {
			vv := *v
			if vv == 0 {
				vv = 7
				*v = vv
			}
			return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
				BitSize:          32,
				Minimum:          func(v int64) *int64 { return &v }(-2147483648),
				Maximum:          func(v int64) *int64 { return &v }(1024),
				ExclusiveMaximum: true,
				ExclusiveMinimum: true,
			})).Validate(int64(vv))
		}(&v.Int), "int")
		errSet.AddErr(func(v **string) error {
			vv := *v
			if vv == nil {
				return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
			}
			return func(v *string) error {
				vv := *v
				if vv == "" {
					return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
				}
				return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
					MinLength: 2,
				})).Validate(string(vv))
			}(vv)
		}(&v.Name), "name")
		errSet.AddErr(func(v ***string) error {
			vv := *v
			if vv == nil {
				var vvv *string
				vv = &vvv
				*v = vv
			}
			return func(v **string) error {
				vv := *v
				if vv == nil {
					var vvv string
					vv = &vvv
					*v = vv
				}
				return func(v *string) error {
					vv := *v
					if vv == "" {
						vv = "11"
						*v = vv
					}
					return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
						Pattern: "\\d+",
					})).Validate(string(vv))
				}(vv)
			}(vv)
		}(&v.ID), "id")
		errSet.AddErr(func(v *github_com_go_courier_schema_testdata_b.PullPolicy) error {
			vv := *v
			if vv == "" {
				return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
			}
			return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
				Enums: []string{
					"Always",
				},
			})).Validate(string(vv))
		}(&v.PullPolicy), "pullPolicy")
		errSet.AddErr(func(v *Protocol) error {
			vv, _ := v.MarshalText()
			if vv != nil {
				return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
					Enums: []string{
						"HTTP",
					},
				})).Validate(string(vv))
			}
			return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
		}(&v.Protocol), "protocol")
		errSet.AddErr(func(v *[]int) error {
			list := *v
			n := len(list)
			if n < 1 {
				return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
					Target:  "slice length",
					Current: n,
					Minimum: 1,
				})
			}
			if n > 30 {
				return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
					Target:  "slice length",
					Current: n,
					Maximum: 30,
				})
			}
			errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()

			validateItem := func(v *int) error {
				vv := *v
				if vv == 0 {
					return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
				}
				return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
					BitSize: 32,
					Minimum: func(v int64) *int64 { return &v }(0),
					Maximum: func(v int64) *int64 { return &v }(10),
				})).Validate(int64(vv))
			}

			for i := 0; i < n; i++ {
				errSet.AddErr(validateItem(&list[i]), i)
			}

			if errSet.Len() > 0 {
				return errSet
			}
			return nil
		}(&v.Slice), "slice")
		errSet.AddErr(func(v *map[string]map[string]struct {
			ID int `json:"id" validate:"@int[0,10]"`
		}) error {
			m := *v
			if n := len(m); n > 3 {
				return &(github_com_go_courier_schema_pkg_validator.OutOfRangeError{
					Target:  "slice length",
					Current: n,
					Maximum: 3,
				})
			}
			errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()

			validateElem := func(v *map[string]struct {
				ID int `json:"id" validate:"@int[0,10]"`
			}) error {
				m := *v
				errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()

				validateElem := func(v *struct {
					ID int `json:"id" validate:"@int[0,10]"`
				}) error {
					errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
					errSet.AddErr(func(v *int) error {
						vv := *v
						if vv == 0 {
							return &(github_com_go_courier_schema_pkg_validator.MissingRequired{})
						}
						return (&(github_com_go_courier_schema_pkg_validator.IntValidator{
							BitSize: 32,
							Minimum: func(v int64) *int64 { return &v }(0),
							Maximum: func(v int64) *int64 { return &v }(10),
						})).Validate(int64(vv))
					}(&v.ID), "id")

					if errSet.Len() > 0 {
						return errSet
					}
					return nil
				}

				for k := range m {

					value := m[k]
					if e := validateElem(&value); e != nil {
						errSet.AddErr(e, k)
					}
				}

				if errSet.Len() > 0 {
					return errSet
				}
				return nil
			}

			for k := range m {

				value := m[k]
				if e := validateElem(&value); e != nil {
					errSet.AddErr(e, k)
				}
			}

			if errSet.Len() > 0 {
				return errSet
			}
			return nil
		}(&v.Map), "map")

		if errSet.Len() > 0 {
			return errSet
		}
		return nil
	}(v)
}
