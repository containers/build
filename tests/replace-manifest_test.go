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
	"encoding/json"
	"os"
	"testing"
)

func TestReplaceManifest(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	manblob, err := json.Marshal(detailedManifest())
	if err != nil {
		panic(err)
	}

	testReplaceWithManifest(t, workingDir, manblob)
}

func TestReplaceManifestWhitespace(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	manblob, err := json.MarshalIndent(detailedManifest(), "", "    ")
	if err != nil {
		panic(err)
	}

	testReplaceWithManifest(t, workingDir, manblob)
}

func testReplaceWithManifest(t *testing.T, workingDir string, manblob []byte) {
	file := mustTempFile()
	defer os.Remove(file.Name())

	_, err := file.Write(manblob)
	if err != nil {
		panic(err)
	}
	file.Close()

	err = runACBuildNoHist(workingDir, "replace-manifest", file.Name())
	if err != nil {
		t.Fatalf("%v", err)
	}

	checkManifest(t, workingDir, detailedManifest())
	checkEmptyRootfs(t, workingDir)
}
