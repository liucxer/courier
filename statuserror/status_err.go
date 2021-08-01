package statuserror

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func IsStatusErr(err error) (*StatusErr, bool) {
	if err == nil {
		return nil, false
	}

	if statusError, ok := err.(StatusError); ok {
		return statusError.StatusErr(), ok
	}

	statusErr, ok := err.(*StatusErr)
	return statusErr, ok
}

func FromErr(err error) *StatusErr {
	if err == nil {
		return nil
	}
	if statusErr, ok := IsStatusErr(err); ok {
		return statusErr
	}
	return Wrap(err, http.StatusInternalServerError, "UnknownError", "unknown error")
}

func Wrap(err error, code int, key string, msgAndDesc ...string) *StatusErr {
	if err == nil {
		return nil
	}

	if len(strconv.Itoa(code)) == 3 {
		code = code * 1e6
	}

	msg := key

	if len(msgAndDesc) > 0 {
		msg = msgAndDesc[0]
	}

	desc := ""

	if len(msgAndDesc) > 1 {
		desc = strings.Join(msgAndDesc[1:], "\n")
	} else {
		desc = err.Error()
	}

	// err = errors.WithMessage(err, "asdfasdfasdfasdfass")

	s := &StatusErr{
		Key:   key,
		Code:  code,
		Msg:   msg,
		Desc:  desc,
		error: errors.WithStack(err),
	}

	return s
}

type StatusErr struct {
	// key of err
	Key string `json:"key" xml:"key"`
	// http code
	Code int `json:"code" xml:"code"`
	// msg of err
	Msg string `json:"msg" xml:"msg"`
	// desc of err
	Desc string `json:"desc" xml:"desc"`
	// can be task error
	// for client to should error msg to end user
	CanBeTalkError bool `json:"canBeTalkError" xml:"canBeTalkError"`

	// request ID or other request context
	ID string `json:"id" xml:"id"`
	// error tracing
	Sources []string `json:"sources" xml:"sources"`
	// error in where fields
	ErrorFields ErrorFields `json:"errorFields" xml:"errorFields"`

	error error
}

func (statusErr *StatusErr) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			e := statusErr.Unwrap()
			if w, ok := e.(WithStackTrace); ok {
				stackTrace := w.StackTrace()
				if len(stackTrace) > 1 {
					_, _ = fmt.Fprintf(s, statusErr.Error()+"%+v", stackTrace[1:])
				}
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, statusErr.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", statusErr.Error())
	}
}

type WithStackTrace interface {
	StackTrace() errors.StackTrace
}

func (statusErr *StatusErr) Unwrap() error {
	return statusErr.error
}

func (statusErr *StatusErr) Summary() string {
	s := fmt.Sprintf(
		`@StatusErr[%s][%d][%s]`,
		statusErr.Key,
		statusErr.Code,
		statusErr.Msg,
	)

	if statusErr.CanBeTalkError {
		return s + "!"
	}
	return s
}

func (statusErr *StatusErr) Is(err error) bool {
	e := FromErr(err)
	if statusErr == nil || e == nil {
		return false
	}
	return e.Key == statusErr.Key && e.Code == statusErr.Code
}

func StatusCodeFromCode(code int) int {
	strCode := fmt.Sprintf("%d", code)
	if len(strCode) < 3 {
		return 0
	}
	statusCode, _ := strconv.Atoi(strCode[:3])
	return statusCode
}

func (statusErr *StatusErr) StatusCode() int {
	return StatusCodeFromCode(statusErr.Code)
}

func (statusErr *StatusErr) Error() string {
	s := fmt.Sprintf(
		"[%s]%s%s",
		strings.Join(statusErr.Sources, ","),
		statusErr.Summary(),
		statusErr.ErrorFields,
	)

	if statusErr.Desc != "" {
		s += " " + statusErr.Desc
	}

	return s
}

func (statusErr StatusErr) WithMsg(msg string) *StatusErr {
	statusErr.Msg = msg
	return &statusErr
}

func (statusErr StatusErr) WithDesc(desc string) *StatusErr {
	statusErr.Desc = desc
	return &statusErr
}

func (statusErr StatusErr) WithID(id string) *StatusErr {
	statusErr.ID = id
	return &statusErr
}

func (statusErr StatusErr) AppendSource(sourceName string) *StatusErr {
	length := len(statusErr.Sources)
	if length == 0 || statusErr.Sources[length-1] != sourceName {
		statusErr.Sources = append(statusErr.Sources, sourceName)
	}
	return &statusErr
}

func (statusErr StatusErr) EnableErrTalk() *StatusErr {
	statusErr.CanBeTalkError = true
	return &statusErr
}

func (statusErr StatusErr) DisableErrTalk() *StatusErr {
	statusErr.CanBeTalkError = false
	return &statusErr
}

func (statusErr StatusErr) AppendErrorField(in string, field string, msg string) *StatusErr {
	statusErr.ErrorFields = append(statusErr.ErrorFields, NewErrorField(in, field, msg))
	return &statusErr
}

func (statusErr StatusErr) AppendErrorFields(errorFields ...*ErrorField) *StatusErr {
	statusErr.ErrorFields = append(statusErr.ErrorFields, errorFields...)
	return &statusErr
}
