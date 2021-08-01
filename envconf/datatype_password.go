package envconf

type Password string

func (p Password) String() string {
	return string(p)
}

func (p Password) SecurityString() string {
	var r []rune
	for range []rune(string(p)) {
		r = append(r, []rune("-")...)
	}
	return string(r)
}
