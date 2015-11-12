// Copyright 2015 The appc Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"os"
	"path"
	"testing"
)

func TestEnd(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	_, _, _, err := runACBuild(workingDir, "end")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	_, err = os.Stat(path.Join(workingDir, ".acbuild"))
	switch {
	case os.IsNotExist(err):
		return
	case err != nil:
		panic(err)
	default:
		t.Fatalf("end failed to remove acbuild working directory")
	}
}
