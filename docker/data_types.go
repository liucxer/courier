package docker

func MaybeListStringFromStringList(list ...string) MaybeListString {
	return MaybeListString{
		values: list,
	}
}

type MaybeListString struct {
	values []string
}

func (v MaybeListString) IsZero() bool {
	return len(v.values) == 0
}

func (v MaybeListString) Value() []string {
	return v.values
}

func (v MaybeListString) MarshalYAML() (interface{}, error) {
	if len(v.values) > 1 {
		return v.values, nil
	}
	return v.values[0], nil
}

func (v *MaybeListString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err == nil {
		v.values = []string{s}
	} else {
		var values []string
		err := unmarshal(&values)
		if err != nil {
			return err
		}
		v.values = values
	}
	return nil
}
