package snapshotmacther

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var snapshotDir = "__snapshots__"
var EnvKeyUpdateSnapshot = "UPDATE_SNAPSHOT"

func MatchSnapshot(names ...string) *SnapshotMatcher {
	filename := filepath.Join(names...)
	updateSnapShot := os.Getenv(EnvKeyUpdateSnapshot)

	return &SnapshotMatcher{
		filename:    filename,
		forceUpdate: updateSnapShot == "all" || updateSnapShot == filename,
	}
}

type SnapshotMatcher struct {
	filename    string
	forceUpdate bool

	actualData   []byte
	expectedData []byte
}

func (matcher *SnapshotMatcher) Match(actual interface{}) (success bool, err error) {
	f := filepath.Join(snapshotDir, matcher.filename)

	switch v := actual.(type) {
	case []byte:
		matcher.actualData = v
	case string:
		matcher.actualData = []byte(v)
	default:
		return false, fmt.Errorf("snapshot not support %T", actual)
	}

	if matcher.forceUpdate {
		if err := snapshot(f, matcher.actualData); err != nil {
			return false, err
		}
		return true, nil
	}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			if err := snapshot(f, matcher.actualData); err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	matcher.expectedData = data

	return bytes.Equal(matcher.actualData, matcher.expectedData), nil
}

func (matcher *SnapshotMatcher) FailureMessage(actual interface{}) (message string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(matcher.expectedData), string(matcher.actualData), true)

	return dmp.DiffPrettyText(diffs) + matcher.helper()
}

func (matcher *SnapshotMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(matcher.expectedData), string(matcher.actualData), true)
	return dmp.DiffPrettyText(diffs) + matcher.helper()
}

func (matcher *SnapshotMatcher) helper() string {
	return fmt.Sprintf(`
-----------
run with "%s=%s" for updating the snapshot or "%s=all" for updating all snapshot  
`, EnvKeyUpdateSnapshot, matcher.filename, EnvKeyUpdateSnapshot)
}

func snapshot(filename string, input []byte) error {
	dirname := filepath.Dir(filename)
	if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, input, os.ModePerm); err != nil {
		return err
	}
	return nil
}
