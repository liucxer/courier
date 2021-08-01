/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"testing"

	"github.com/onsi/gomega"
)

func TestExtractCommentTags(t *testing.T) {
	commentLines := []string{
		"Human comment that is ignored.",
		"+gengo:test=value1",
		"+bar",
		"+baz=qux,zrb=true",
		"+gengo:test=value2",
	}

	a, others := ExtractCommentTags("+", commentLines)

	gomega.NewWithT(t).Expect(a).To(gomega.Equal(map[string][]string{
		"gengo:test": {"value1", "value2"},
		"bar":        {""},
		"baz":        {"qux,zrb=true"},
	}))

	gomega.NewWithT(t).Expect(others).To(gomega.Equal([]string{
		"Human comment that is ignored.",
	}))
}
