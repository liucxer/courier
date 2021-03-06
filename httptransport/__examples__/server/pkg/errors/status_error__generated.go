package errors

import (
	github_com_go_courier_statuserror "github.com/liucxer/courier/statuserror"
)

var _ interface {
	github_com_go_courier_statuserror.StatusError
} = (*StatusError)(nil)

func (v StatusError) StatusErr() *github_com_go_courier_statuserror.StatusErr {
	return &github_com_go_courier_statuserror.StatusErr{
		Key:            v.Key(),
		Code:           v.Code(),
		Msg:            v.Msg(),
		CanBeTalkError: v.CanBeTalkError(),
	}
}

func (v StatusError) Error() string {
	return v.StatusErr().Error()
}

func (v StatusError) Unwrap() error {
	return v.StatusErr()
}

func (v StatusError) StatusCode() int {
	return github_com_go_courier_statuserror.StatusCodeFromCode(int(v))
}

func (v StatusError) Code() int {
	if withServiceCode, ok := (interface{})(v).(github_com_go_courier_statuserror.StatusErrorWithServiceCode); ok {
		return withServiceCode.ServiceCode() + int(v)
	}
	return int(v)

}

func (v StatusError) Key() string {
	switch v {
	case Unauthorized:
		return "Unauthorized"
	case InternalServerError:
		return "InternalServerError"
	}
	return "UNKNOWN"
}

func (v StatusError) Msg() string {
	switch v {
	case Unauthorized:
		return "Unauthorized"
	case InternalServerError:
		return "InternalServerError"
	}
	return "-"
}

func (v StatusError) CanBeTalkError() bool {
	switch v {
	case Unauthorized:
		return true
	case InternalServerError:
		return false
	}
	return false
}
