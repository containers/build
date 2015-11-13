// Copyright 2015 Simone Gotti
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

type ChangeType uint8

const (
	Added ChangeType = iota
	Modified
	Deleted
)

// FSChanges represents a slice of changes
// This isn't a map to keep ordering on the changes, for a map see FSChangesMap.
type FSChanges []*FSChange

type FSChange struct {
	Path string
	ChangeType
}

// FSChanges represents a map of changes, the map's key is the path while the
// value is the ChangeType
type FSChangesMap map[string]ChangeType

// Utility function to convert an FSChanges slice to a FSChangesMap
func (fsc FSChanges) ToMap() FSChangesMap {
	fscm := make(FSChangesMap, len(fsc))
	for _, c := range fsc {
		fscm[c.Path] = c.ChangeType
	}
	return fscm
}

// The FSDiffer interface should be implemented from an fsdiffer implementation
// The returned FSChanges should be lexically ordered like filepath.Walk() does.
type FSDiffer interface {
	Diff() (FSChanges, error)
}
