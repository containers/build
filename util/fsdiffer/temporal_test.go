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
package fsdiffer

import (
	"archive/tar"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// Do not test (at the moment) symlinks as every os needs a special syscall to set symlink times
func TestTemporalFSDiffer(t *testing.T) {
	time1 := time.Now()
	time2 := time.Now().Add(1 * time.Microsecond)
	sourceFiles := []*buildFileInfo{
		&buildFileInfo{path: "file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
		&buildFileInfo{path: "dir01", typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time1},
		&buildFileInfo{path: "dir01/file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
	}

	tests := []struct {
		sourceFiles     []*buildFileInfo
		destFiles       []*buildFileInfo
		expectedChanges FSChangesMap
	}{
		{
			sourceFiles: sourceFiles,
			destFiles: []*buildFileInfo{
				&buildFileInfo{path: "file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
				&buildFileInfo{path: "dir01", typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time1},
				&buildFileInfo{path: "dir01/file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
			},
			expectedChanges: FSChangesMap{},
		},
		{
			sourceFiles: sourceFiles,
			// dir01 mtime changed
			destFiles: []*buildFileInfo{
				&buildFileInfo{path: "file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
				&buildFileInfo{path: "dir01", typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time2},
				&buildFileInfo{path: "dir01/file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
			},
			expectedChanges: FSChangesMap{"dir01": Modified},
		},
		{
			sourceFiles: sourceFiles,

			// file01 contents changed
			destFiles: []*buildFileInfo{
				&buildFileInfo{path: "file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hellohello"},
				&buildFileInfo{path: "dir01", typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time1},
				&buildFileInfo{path: "dir01/file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
			},
			expectedChanges: FSChangesMap{"file01": Modified},
		},
		{
			sourceFiles: sourceFiles,
			// new dir and file dir02/file01, dir01/file01 removed
			destFiles: []*buildFileInfo{
				&buildFileInfo{path: "file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
				&buildFileInfo{path: "dir02", typeflag: tar.TypeDir, mode: 0755, atime: time1, mtime: time1},
				&buildFileInfo{path: "dir02/file01", typeflag: tar.TypeReg, mode: 0644, atime: time1, mtime: time1, contents: "hello"},
			},
			expectedChanges: FSChangesMap{
				"dir01":        Deleted,
				"dir01/file01": Deleted,
				"dir02":        Added,
				"dir02/file01": Added,
			},
		},
	}
	dir, err := ioutil.TempDir("", tstprefix)
	if err != nil {
		t.Fatalf("error creating tempdir: %v", err)
	}
	defer os.RemoveAll(dir)
	testDir := filepath.Join(dir, "test")

	for _, tt := range tests {
		os.RemoveAll(testDir)

		err = buildFS(testDir, tt.sourceFiles)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		tm, err := NewTemporalFSDiffer(testDir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		os.RemoveAll(testDir)
		err = buildFS(testDir, tt.destFiles)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		changes, err := tm.Diff()
		changesMap := changes.ToMap()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(changesMap) != len(tt.expectedChanges) {
			t.Errorf("wrong changes size: want: %d, got: %d", len(tt.expectedChanges), len(changesMap))
		}
		if !reflect.DeepEqual(changesMap, tt.expectedChanges) {
			t.Errorf("changes differs: want: %q, got: %q", printChanges(tt.expectedChanges), printChanges(changesMap))
		}
	}

}
