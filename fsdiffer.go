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
