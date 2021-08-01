package statuserror

type StatusError interface {
	StatusErr() *StatusErr
	Error() string
}

type StatusErrorWithServiceCode interface {
	ServiceCode() int
}
