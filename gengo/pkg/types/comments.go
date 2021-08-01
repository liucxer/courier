package types

import (
	"strings"
)

// ExtractCommentTags parses comments for lines of the form:
//
//   'marker' + "key=value".
//
// Values are optional; "" is the default.  A tag can be specified more than
// one time and all values are returned.  If the resulting map has an entry for
// a key, the value (a slice) is guaranteed to have at least 1 element.
//
// Example: if you pass "+" for 'marker', and the following lines are in
// the comments:
//   +foo=value1
//   +bar
//   +foo=value2
//   +baz="qux"
// Then this function will return:
//   map[string][]string{"foo":{"value1, "value2"}, "bar": {"true"}, "baz": {"qux"}}
func ExtractCommentTags(marker string, lines []string) (map[string][]string, []string) {
	out := map[string][]string{}
	otherLines := make([]string, 0)

	for _, line := range lines {
		line = strings.Trim(line, " ")

		if !strings.HasPrefix(line, marker) {
			otherLines = append(otherLines, line)
			continue
		}

		kv := strings.SplitN(line[len(marker):], "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = append(out[kv[0]], kv[1])
		} else if len(kv) == 1 {
			out[kv[0]] = append(out[kv[0]], "")
		}
	}

	return out, otherLines
}
