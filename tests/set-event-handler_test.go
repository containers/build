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

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

func manWithEH(eh types.EventHandler) schema.ImageManifest {
	man := emptyManifestWithApp()
	man.App.EventHandlers = []types.EventHandler{eh}
	return man
}

func TestEHPreStart(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	var prestart = []string{"/bin/sh", "-c", "'chmod 755 /'"}

	err := runACBuildNoHist(workingDir, append([]string{"set-event-handler", "pre-start", "--"}, prestart...)...)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	man := manWithEH(types.EventHandler{
		Name: "pre-start",
		Exec: prestart,
	})

	checkManifest(t, workingDir, man)
	checkEmptyRootfs(t, workingDir)
}

func TestEHPostStop(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	var poststop = []string{"/bin/rm", "-rf", "/Plex Media Server/plexmediaserver.pid"}

	err := runACBuildNoHist(workingDir, append([]string{"set-event-handler", "post-stop", "--"}, poststop...)...)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	man := manWithEH(types.EventHandler{
		Name: "post-stop",
		Exec: poststop,
	})

	checkManifest(t, workingDir, man)
	checkEmptyRootfs(t, workingDir)
}
