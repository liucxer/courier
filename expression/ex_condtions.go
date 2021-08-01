package expression

import (
	"fmt"
	"regexp"
)

func init() {
	Register("in", func(f ExFactory, args []interface{}) (ex ExDo, e error) {
		return BuildEx(`["in", target, any, any, ...]`, 2, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				target, err := exes[0](params...)
				if err != nil {
					return nil, err
				}
				values, err := Exec(exes[1:], params...)
				if err != nil {
					return nil, err
				}
				for i := range values {
					if values[i] == target {
						return true, nil
					}
				}
				return false, nil
			}, nil
		})(f, args)
	})

	Register("match", func(f ExFactory, args []interface{}) (ex ExDo, e error) {
		return BuildEx(`["match", "regexp", v]`, 2, func(exes []ExDo) (ex ExDo, e error) {
			arg, err := exes[0]()
			if err != nil {
				return nil, err
			}

			reg, ok := arg.(string)
			if !ok {
				return nil, NewInvalidExpression("regexp must be string")
			}

			r, err := regexp.Compile(reg)
			if err != nil {
				return nil, NewInvalidExpression(err.Error())
			}

			return func(params ...interface{}) (interface{}, error) {
				arg, err := exes[1](params...)
				if err != nil {
					return nil, err
				}

				s, err := toString(arg)
				if err != nil {
					return nil, err
				}

				return r.MatchString(s), nil
			}, nil
		})(f, args)
	})

	Register("allOf", func(f ExFactory, args []interface{}) (ex ExDo, e error) {
		return BuildEx(`["allOf", any, any, ...]`, 1, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				conditions := make([]Condition, len(exes))

				for i := range conditions {
					ex := exes[i]
					conditions[i] = func() (b bool, e error) {
						v, err := ex(params...)
						if err != nil {
							return false, err
						}
						return toBoolean(v)
					}
				}

				return AllOf(conditions...)
			}, nil
		})(f, args)
	})

	Register("anyOf", func(f ExFactory, args []interface{}) (ex ExDo, e error) {
		return BuildEx(`["anyOf", any, any, ...]`, 1, func(exes []ExDo) (ex ExDo, e error) {
			return func(params ...interface{}) (interface{}, error) {
				conditions := make([]Condition, len(exes))

				for i := range conditions {
					ex := exes[i]
					conditions[i] = func() (b bool, e error) {
						v, err := ex(params...)
						if err != nil {
							return false, err
						}
						return toBoolean(v)
					}
				}

				return AnyOf(conditions...)
			}, nil
		})(f, args)
	})

	Register("not", createConvert(func(v interface{}) (interface{}, error) {
		b, err := toBoolean(v)
		return !b, err
	}))

	Register("eq", createExCompare("eq", func(left interface{}, right interface{}) (bool, error) {
		return AnyOf(
			func() (b bool, e error) {
				return compareNumberic(left, right, func(left float64, right float64) bool {
					return left == right
				})
			},
			func() (b bool, e error) {
				return compareStringable(left, right, func(left string, right string) bool {
					return left == right
				})
			},
		)
	}))

	Register("gt", createExCompare("gt", func(left interface{}, right interface{}) (bool, error) {
		return compareNumberic(left, right, func(left float64, right float64) bool {
			return left > right
		})
	}))

	Register("gte", createExCompare("gte", func(left interface{}, right interface{}) (bool, error) {
		return compareNumberic(left, right, func(left float64, right float64) bool {
			return left >= right
		})
	}))

	Register("lt", createExCompare("lt", func(left interface{}, right interface{}) (bool, error) {
		return compareNumberic(left, right, func(left float64, right float64) bool {
			return left < right
		})
	}))

	Register("lte", createExCompare("lte", func(left interface{}, right interface{}) (bool, error) {
		return compareNumberic(left, right, func(left float64, right float64) bool {
			return left <= right
		})
	}))

	Register("neq", createExCompare("neq", func(left interface{}, right interface{}) (bool, error) {
		return compareNumberic(left, right, func(left float64, right float64) bool {
			return left != right
		})
	}))
}
func createExCompare(key string, compare func(left interface{}, right interface{}) (bool, error)) ExBuilder {
	return BuildEx(fmt.Sprintf(`["%s", a: number, b: number]`, key), 2, func(exes []ExDo) (ex ExDo, e error) {
		return func(in ...interface{}) (interface{}, error) {
			values, err := Exec(exes, in...)
			if err != nil {
				return nil, err
			}
			return compare(values[0], values[1])
		}, nil
	})
}

func compareNumberic(left interface{}, right interface{}, compare func(left float64, right float64) bool) (bool, error) {
	leftN, err := toNumber(left)
	if err != nil {
		return false, err
	}
	rightN, err := toNumber(right)
	if err != nil {
		return false, err
	}
	return compare(leftN, rightN), nil
}

func compareStringable(left interface{}, right interface{}, compare func(left string, right string) bool) (bool, error) {
	leftS, err := toString(left)
	if err != nil {
		return false, err
	}
	rights, err := toString(right)
	if err != nil {
		return false, err
	}
	return compare(leftS, rights), nil
}

type Condition func() (bool, error)

func AllOf(conditions ...Condition) (bool, error) {
	for _, cond := range conditions {
		ret, err := cond()
		if err != nil {
			return false, err
		}
		if !ret {
			return false, nil
		}
	}
	return true, nil
}

func AnyOf(conditions ...Condition) (bool, error) {
	for _, cond := range conditions {
		ret, err := cond()
		if err != nil {
			return false, err
		}
		if ret {
			return true, nil
		}
	}
	return false, nil
}
