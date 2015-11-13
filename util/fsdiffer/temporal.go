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
	"path/filepath"
)

// TemporalFSDiffer is used to generate changes in a given directory
// between two different points in time.
type TemporalFSDiffer struct {
	dir    string
	before map[string]fileInfo
}

// NewTemporalFSDiffer creates a new TemporalFSDiffer that will report
// changes on the given directory.
func NewTemporalFSDiffer(dir string) (*TemporalFSDiffer, error) {
	t := &TemporalFSDiffer{
		dir:    dir,
		before: make(map[string]fileInfo),
	}
	err := filepath.Walk(t.dir, fsWalker(t.before))
	if err != nil {
		return nil, err
	}
	return t, nil
}

// Diff will return any changes to the filesystem in the provided directory
// since Start was called.
//
// To detect if a file was changed it checks the file's size and mtime (like
// rsync does by default if no --checksum options is used)
func (t *TemporalFSDiffer) Diff() (FSChanges, error) {
	changes := FSChanges{}
	after := make(map[string]fileInfo)
	err := filepath.Walk(t.dir, fsWalker(after))
	if err != nil {
		return nil, err
	}

	for _, afterInfo := range after {
		relpath, _ := filepath.Rel(t.dir, afterInfo.Path)
		sourceInfo, ok := t.before[filepath.Join(t.dir, relpath)]
		if !ok {
			changes = append(changes, &FSChange{Path: relpath, ChangeType: Added})
		} else {
			if sourceInfo.Size() != afterInfo.Size() || sourceInfo.ModTime().Before(afterInfo.ModTime()) {
				changes = append(changes, &FSChange{Path: relpath, ChangeType: Modified})
			}
		}
	}
	for _, beforeInfo := range t.before {
		relpath, _ := filepath.Rel(t.dir, beforeInfo.Path)
		_, ok := after[filepath.Join(t.dir, relpath)]
		if !ok {
			changes = append(changes, &FSChange{Path: relpath, ChangeType: Deleted})
		}
	}
	return changes, nil
}
