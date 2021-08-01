package tag

type Tag struct {
	Key   string `db:"f_key" json:"key"`
	Value string `db:"f_value" json:"value"`
}

type Tags map[string][]string

func (Tags) New() interface{} {
	return &Tag{}
}

func (u Tags) Get(k string) string {
	if values, ok := u[k]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (u Tags) Add(k, v string) {
	if _, ok := u[k]; !ok {
		u[k] = []string{}
	}
	u[k] = append(u[k], v)
}

func (u Tags) Set(k string, values ...string) {
	u[k] = values
}

func (u Tags) Next(v interface{}) error {
	tag := v.(*Tag)
	u.Add(tag.Key, tag.Value)
	return nil
}
