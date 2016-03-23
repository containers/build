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
	"io/ioutil"
	"path"
	"testing"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/schema"
)

func TestCatEmptyManifest(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	testCat(t, workingDir)
}

func TestCatDetailedManifest(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	manblob, err := json.Marshal(detailedManifest())
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(workingDir, ".acbuild", "currentaci", aci.ManifestFile), manblob, 0644)
	if err != nil {
		panic(err)
	}

	testCat(t, workingDir)
}

func testCat(t *testing.T, workingDir string) {
	wantedManblob, err := ioutil.ReadFile(path.Join(workingDir, ".acbuild", "currentaci", aci.ManifestFile))
	if err != nil {
		panic(err)
	}
	wantedManblob = append(wantedManblob, byte('\n'))

	var man schema.ImageManifest
	err = man.UnmarshalJSON(wantedManblob)
	if err != nil {
		panic(err)
	}

	_, manblob, _, err := runACBuild(workingDir, "cat-manifest")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if manblob != string(wantedManblob) {
		t.Fatalf("printed manifest and manifest on disk differ")
	}

	wantedManblob = prettyPrintMan(wantedManblob)
	wantedManblob = append(wantedManblob, byte('\n'))

	_, manblob, _, err = runACBuild(workingDir, "cat-manifest", "--pretty-print")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if manblob != string(wantedManblob) {
		t.Fatalf("pretty printed manifest and manifest on disk differ")
	}

	checkManifest(t, workingDir, man)
	checkEmptyRootfs(t, workingDir)
}

func prettyPrintMan(manblob []byte) []byte {
	var man schema.ImageManifest

	err := man.UnmarshalJSON(manblob)
	if err != nil {
		panic(err)
	}

	manblob2, err := json.MarshalIndent(man, "", "    ")
	if err != nil {
		panic(err)
	}
	return manblob2
}
