package examples

import (
	"errors"
	"fmt"
	"testing"

	"github.com/liucxer/courier/statuserror"
)

var ErrorNotFound = errors.New("404 NotFound")

func sentinel() error {
	err := ErrorNotFound
	sErr := statuserror.Wrap(err, 400, "httpclient error")
	return sErr
}

func biz() error {
	return sentinel()
}

func api() error {
	return biz()
}

func Test_StatusError(t *testing.T) {
	err := api()

	if errors.Is(err, ErrorNotFound) {
		fmt.Printf("%+v", err)
	}

	// === RUN   Test_StatusError
	// []@StatusErr[httpclient error][400000000][httpclient error] 404 NotFound
	// github.com/liucxer/courier/statuserror/__examples__.sentinel
	// 	/private/tmp/statuserror/__examples__/status_error_wrap_test.go:15
	// github.com/liucxer/courier/statuserror/__examples__.biz
	// 	/private/tmp/statuserror/__examples__/status_error_wrap_test.go:20
	// ....
}
