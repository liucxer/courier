package b

import (
	github_com_go_courier_schema_pkg_validator "github.com/liucxer/courier/schema/pkg/validator"
)

func (v *GetByID) Validate() error {
	return func(v *GetByID) error {
		errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()
		errSet.AddErr(func(v *string) error {
			vv := *v
			if vv == "" {
				return nil
			}
			return (&(github_com_go_courier_schema_pkg_validator.StringValidator{
				MinLength: 4,
			})).Validate(string(vv))
		}(&v.Name), github_com_go_courier_schema_pkg_validator.Location("query"), "name")
		errSet.AddErr(func(v *[]string) error {
			list := *v
			n := len(list)
			errSet := github_com_go_courier_schema_pkg_validator.NewErrorSet()

			validateItem := func(v *string) error {
				return nil
			}

			for i := 0; i < n; i++ {
				errSet.AddErr(validateItem(&list[i]), i)
			}

			if errSet.Len() > 0 {
				return errSet
			}
			return nil
		}(&v.Label), github_com_go_courier_schema_pkg_validator.Location("query"), "label")

		if errSet.Len() > 0 {
			return errSet
		}
		return nil
	}(v)
}
