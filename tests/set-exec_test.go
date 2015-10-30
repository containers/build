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
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

func TestSetExec(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	var exec = []string{"/bin/nethack4", "-D", "wizard"}

	err := runACBuild(workingDir, append([]string{"set-exec", "--"}, exec...)...)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	man := schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:  exec,
			User:  "0",
			Group: "0",
		},
		Labels: systemLabels,
	}

	checkManifest(t, workingDir, man)
	checkEmptyRootfs(t, workingDir)
}
