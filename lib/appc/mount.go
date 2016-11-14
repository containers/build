// Copyright 2016 The appc Authors
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

package appc

import (
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

// AddMount will add a mount point with the given name and path to the untarred
// ACI stored at a.CurrentImagePath. If the mount point already exists its value
// will be updated to the new value. readOnly signifies whether or not the
// mount point should be read only.
func (m *Manifest) AddMount(name, path string, readOnly bool) error {
	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}

	removeMount(name, m.manifest)
	removeMount(path, m.manifest)
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	m.manifest.App.MountPoints = append(m.manifest.App.MountPoints,
		types.MountPoint{
			Name:     *acn,
			Path:     path,
			ReadOnly: readOnly,
		})
	return m.save()
}

// RemoveMount will remove the mount point with the given name from the
// untarred ACI stored at a.CurrentImagePath
func (m *Manifest) RemoveMount(mount string) error {
	err := removeMount(mount, m.manifest)
	if err != nil {
		return err
	}

	return m.save()
}

func removeMount(mount string, m *schema.ImageManifest) error {
	if m.App == nil {
		return ErrNotFound
	}

	foundOne := false
	for i := len(m.App.MountPoints) - 1; i >= 0; i-- {
		if string(m.App.MountPoints[i].Name) == mount || m.App.MountPoints[i].Path == mount {
			foundOne = true
			m.App.MountPoints = append(
				m.App.MountPoints[:i],
				m.App.MountPoints[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
