package statuserror

import (
	"fmt"
	"regexp"
	"strconv"
)

func ParseStatusErrSummary(s string) (*StatusErr, error) {
	if !reStatusErrSummary.Match([]byte(s)) {
		return nil, fmt.Errorf("unsupported status err summary: %s", s)
	}

	matched := reStatusErrSummary.FindStringSubmatch(s)

	code, _ := strconv.ParseInt(matched[2], 10, 64)

	return &StatusErr{
		Key:            matched[1],
		Code:           int(code),
		Msg:            matched[3],
		CanBeTalkError: matched[4] != "",
	}, nil
}

// @err[UnknownError][500000000][unknown error]
var reStatusErrSummary = regexp.MustCompile(`@StatusErr\[(.+)\]\[(.+)\]\[(.+)\](!)?`)