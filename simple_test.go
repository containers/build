package fsdiffer

import (
	"archive/tar"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

const tstprefix = "simplefsdiffer-test"

type buildFileInfo struct {
	path     string
	typeflag byte
	size     int64
	mode     os.FileMode
	atime    time.Time
	mtime    time.Time
	contents string
}

func buildFS(dir string, files []*buildFileInfo) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	dirs := []*buildFileInfo{}
	for _, f := range files {
		p := filepath.Join(dir, f.path)
		switch f.typeflag {
		case tar.TypeDir:
			dirs = append(dirs, f)
			if err := os.MkdirAll(p, f.mode); err != nil {
				return err
			}
			dir, err := os.Open(p)
			if err != nil {
				return err
			}
			if err := dir.Chmod(f.mode); err != nil {
				dir.Close()
				return err
			}
			dir.Close()
		case tar.TypeReg:
			err := ioutil.WriteFile(p, []byte(f.contents), f.mode)
			if err != nil {
				return err
			}

		}

		if err := os.Chtimes(p, f.atime, f.mtime); err != nil {
			return err
		}
	}

	// Restore dirs atime and mtime. This has to be done after extracting
	// as a file extraction will change its parent directory's times.
	for _, d := range dirs {
		p := filepath.Join(dir, d.path)
		if err := os.Chtimes(p, d.atime, d.mtime); err != nil {
			return err
		}
	}

	defaultTime := time.Unix(0, 0)
	// Restore the base dir time as it will be changed by the previous extractions
	if err := os.Chtimes(dir, defaultTime, defaultTime); err != nil {
		return err
	}
	return nil
}

func fixeDirTimes(dir string) {

}

// Do not test (at the moment) symlinks as every os needs a special syscall to set symlink times
func TestSimpleFSDiffer(t *testing.T) {
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
	sourceDir := filepath.Join(dir, "source")
	destDir := filepath.Join(dir, "dest")

	for _, tt := range tests {
		err = os.RemoveAll(sourceDir)
		err = os.RemoveAll(destDir)

		err = buildFS(sourceDir, tt.sourceFiles)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = buildFS(destDir, tt.destFiles)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		s := NewSimpleFSDiffer(sourceDir, destDir)
		changes, err := s.Diff()
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

func printChanges(changes FSChangesMap) string {
	changesStr := make([]string, len(changes))
	i := 0
	for p, c := range changes {
		changesStr[i] = fmt.Sprintf("%s: %d ", p, c)
		i++
	}
	return strings.Join(changesStr, ", ")
}
