package a

// +gengo:enum
type Protocol int

func (Protocol) Offset() int {
	return -4
}

const (
	PROTOCOL_UNKNOWN Protocol = iota
	PROTOCOL__HTTP            // http
	PROTOCOL__HTTPS           // https
	_
	_1
	_2
	PROTOCOL__TCP // tcp
)
