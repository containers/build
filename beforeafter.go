package fsdiffer

import (
	"fmt"
	"path/filepath"
)

// The BeforeAfterFSDiffer is used to generate changes in a given directory
// between two different points in time.
type BeforeAfterFSDiffer struct {
	dir    string
	before map[string]fileInfo
}

// NewBeforeAfterFSDiffer creates a new BeforeAfterFSDiffer that will report
// changes on the given directory.
func NewBeforeAfterFSDiffer(dir string) *BeforeAfterFSDiffer {
	return &BeforeAfterFSDiffer{dir: dir}
}

// Start will walk the provided directory, and record its state for use as the
// before part of the changes returned by Diff.
func (ba *BeforeAfterFSDiffer) Start() error {
	if ba.before != nil {
		return fmt.Errorf("start already called")
	}
	ba.before = make(map[string]fileInfo)
	return filepath.Walk(ba.dir, fsWalker(ba.before))
}

// Diff will return any changes to the filesystem in the provided directory
// since Start was called.
//
// To detect if a file was changed it checks the file's size and mtime (like
// rsync does by default if no --checksum options is used)
func (ba *BeforeAfterFSDiffer) Diff() (FSChanges, error) {
	changes := FSChanges{}
	after := make(map[string]fileInfo)
	err := filepath.Walk(ba.dir, fsWalker(after))
	if err != nil {
		return nil, err
	}

	for _, afterInfo := range after {
		relpath, _ := filepath.Rel(ba.dir, afterInfo.Path)
		sourceInfo, ok := ba.before[filepath.Join(ba.dir, relpath)]
		if !ok {
			changes = append(changes, &FSChange{Path: relpath, ChangeType: Added})
		} else {
			if sourceInfo.Size() != afterInfo.Size() || sourceInfo.ModTime().Before(afterInfo.ModTime()) {
				changes = append(changes, &FSChange{Path: relpath, ChangeType: Modified})
			}
		}
	}
	for _, infoA := range ba.before {
		relpath, _ := filepath.Rel(ba.dir, infoA.Path)
		_, ok := after[filepath.Join(ba.dir, relpath)]
		if !ok {
			changes = append(changes, &FSChange{Path: relpath, ChangeType: Deleted})
		}
	}
	return changes, nil
}
