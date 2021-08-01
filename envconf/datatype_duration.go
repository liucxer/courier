package envconf

import (
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	dur, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}
	*d = Duration(dur)
	return nil
}

func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(time.Duration(d).String()), nil
}
