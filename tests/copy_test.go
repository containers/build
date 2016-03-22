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
	"archive/tar"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"

	"github.com/appc/acbuild/util/fsdiffer"
)

const dest = "/test"

type buildFileInfo struct {
	path     string
	typeflag byte
	size     int64
	mode     os.FileMode
	atime    time.Time
	mtime    time.Time
	contents string
}

func mkBuildFileInfoDir(path string, time1 time.Time) *buildFileInfo {
	return &buildFileInfo{path: path, typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time1}
}

func mkBuildFileInfoFile(path string, time1 time.Time) *buildFileInfo {
	return &buildFileInfo{path: path, typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"}
}

func mustBuildFS(dir string, files []*buildFileInfo) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	dirs := []*buildFileInfo{}
	for _, f := range files {
		p := filepath.Join(dir, f.path)
		switch f.typeflag {
		case tar.TypeDir:
			dirs = append(dirs, f)
			if err := os.MkdirAll(p, f.mode); err != nil {
				panic(err)
			}
			dir, err := os.Open(p)
			if err != nil {
				panic(err)
			}
			if err := dir.Chmod(f.mode); err != nil {
				dir.Close()
				panic(err)
			}
			dir.Close()
		case tar.TypeReg:
			err := ioutil.WriteFile(p, []byte(f.contents), f.mode)
			if err != nil {
				panic(err)
			}

		}

		if err := os.Chtimes(p, f.atime, f.mtime); err != nil {
			panic(err)
		}
	}

	// Restore dirs atime and mtime. This has to be done after extracting
	// as a file extraction will change its parent directory's times.
	for _, d := range dirs {
		p := filepath.Join(dir, d.path)
		if err := os.Chtimes(p, d.atime, d.mtime); err != nil {
			panic(err)
		}
	}

	defaultTime := time.Unix(0, 0)
	// Restore the base dir time as it will be changed by the previous extractions
	if err := os.Chtimes(dir, defaultTime, defaultTime); err != nil {
		panic(err)
	}
}

func TestCopyOneDir(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	sourceDir, err := ioutil.TempDir("", "acbuild-test-copy")
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

	err = runACBuildNoHist(workingDir, "--debug", "copy", sourceDir, dest)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	testMatchingFSTree(t, workingDir, sourceDir, dest)
}

func TestCopyOneFile(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	sourceFile, err := ioutil.TempFile("", "acbuild-test-copy")
	if err != nil {
		panic(err)
	}
	_, err = sourceFile.Write([]byte("this is a test file"))
	if err != nil {
		panic(err)
	}
	sourceFile.Close()
	defer os.RemoveAll(sourceFile.Name())

	err = runACBuildNoHist(workingDir, "copy", sourceFile.Name(), dest)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	testMatchingFSTree(t, workingDir, sourceFile.Name(), dest)
}

func TestCopyOneDirToExistingDir(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	sourceDir, err := ioutil.TempDir("", "acbuild-test-copy")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(sourceDir)

	err = os.MkdirAll(path.Join(workingDir, ".acbuild", "currentaci", aci.RootfsDir, dest), 0755)
	if err != nil {
		panic(err)
	}

	const expectedErrorMsg = "copy: mkdir .acbuild/currentaci/rootfs/test: file exists\n"

	_, _, stderr, err := runACBuild(workingDir, "--no-history", "copy", sourceDir, dest)
	if err != nil && stderr != expectedErrorMsg {
		t.Fatalf("was expecting an error, but not this one:\n%v", err)
	}
	if err == nil {
		t.Fatalf("got no error when destination directory exists, was expecting one")
	}
}

func TestCopyManyDirs(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	sourceDir, err := ioutil.TempDir("", "acbuild-test-copy")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(sourceDir)

	time1 := time.Now()
	files := []*buildFileInfo{
		mkBuildFileInfoDir("dir01", time1),
		mkBuildFileInfoFile("dir01/file01", time1),
		mkBuildFileInfoDir("dir02", time1),
		mkBuildFileInfoFile("dir02/file02", time1),
		mkBuildFileInfoDir("dir02/dir02.1", time1),
		mkBuildFileInfoFile("dir02/dir02.1/file02", time1),
		mkBuildFileInfoDir("dir03", time1),
		mkBuildFileInfoFile("dir03/file03", time1),
	}

	mustBuildFS(sourceDir, files)

	froms := []string{"dir01", "dir02", "dir03"}
	for i := 0; i < len(froms); i++ {
		froms[i] = path.Join(sourceDir, froms[i])
	}

	// golang--
	err = runACBuildNoHist(workingDir, append([]string{"copy-to-dir"}, append(froms, dest)...)...)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	testMatchingFSTree(t, workingDir, sourceDir, dest)
}

func TestCopyDirsAndFiles(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	sourceDir, err := ioutil.TempDir("", "acbuild-test-copy")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(sourceDir)

	time1 := time.Now()
	files := []*buildFileInfo{
		mkBuildFileInfoDir("dir01", time1),
		mkBuildFileInfoDir("dir02", time1),
		mkBuildFileInfoFile("dir02/file02", time1),
		mkBuildFileInfoFile("file01", time1),
		mkBuildFileInfoFile("file02", time1),
	}

	mustBuildFS(sourceDir, files)

	froms := []string{"dir01", "dir02", "file01", "file02"}
	for i := 0; i < len(froms); i++ {
		froms[i] = path.Join(sourceDir, froms[i])
	}

	// golang--
	err = runACBuildNoHist(workingDir, append([]string{"copy-to-dir"}, append(froms, dest)...)...)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	testMatchingFSTree(t, workingDir, sourceDir, dest)
}

func testMatchingFSTree(t *testing.T, workingDir, source, dest string) {
	dest = path.Join(workingDir, ".acbuild", "currentaci", aci.RootfsDir, dest)

	changes, err := fsdiffer.NewSimpleFSDiffer(source, dest).Diff()
	if err != nil {
		t.Fatalf("%v", err)
	}

	// Depending on how acbuild copy was called, the given directory may have a
	// different modification time and/or mode. If the only change was the
	// current directory being modified, let's ignore it.
	if len(changes) == 1 && changes[0].Path == "." && changes[0].ChangeType == fsdiffer.Modified {
		return
	}

	if len(changes) != 0 {
		var changestring string
		for _, change := range changes {
			if changestring != "" {
				changestring += "\n"
			}
			changestring += fmt.Sprintf("file changed: %s, change type: %d", change.Path, change.ChangeType)
		}
		t.Fatalf("Got %d changes, expected 0\n%s", len(changes), changestring)
	}
}
