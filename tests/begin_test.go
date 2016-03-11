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
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestBeginEmpty(t *testing.T) {
	workingDir := mustTempDir()
	defer cleanUpTest(workingDir)

	_, _, _, err := runACBuild(workingDir, "begin")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifest())
	checkEmptyRootfs(t, workingDir)
}

func TestBeginLocalACI(t *testing.T) {
	tmpaci, err := ioutil.TempFile("", "acbuild-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpaci.Name())

	makeACI(tmpaci, detailedManifest())

	tmpaci.Close()

	workingDir := mustTempDir()
	defer cleanUpTest(workingDir)

	_, _, _, err = runACBuild(workingDir, "begin", tmpaci.Name())
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	checkManifest(t, workingDir, detailedManifest())
	checkEmptyRootfs(t, workingDir)
}

func TestBeginLocalDirectory(t *testing.T) {
	sourceDir, err := ioutil.TempDir("", "acbuild-test-begin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(sourceDir)

	time1 := time.Now()
	files := []*buildFileInfo{
		mkBuildFileInfoFile("file01", time1),
		mkBuildFileInfoDir("dir01", time1),
		mkBuildFileInfoFile("dir01/file01", time1),
	}

	mustBuildFS(sourceDir, files)

	workingDir := mustTempDir()
	defer cleanUpTest(workingDir)

	_, _, _, err = runACBuild(workingDir, "begin", sourceDir)
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}

	checkManifest(t, workingDir, emptyManifest())
	testMatchingFSTree(t, workingDir, sourceDir, "/")
}
