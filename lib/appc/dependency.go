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

// AddDependency will add a dependency with the given name, id, labels, and size
// to the untarred ACI stored at a.CurrentImagePath. If the dependency already
// exists its fields will be updated to the new values.
func (m *Manifest) AddDependency(imageName types.ACIdentifier, imageId *types.Hash, labels types.Labels, size uint) error {
	removeDep(imageName, m.manifest)
	m.manifest.Dependencies = append(m.manifest.Dependencies,
		types.Dependency{
			ImageName: imageName,
			ImageID:   imageId,
			Labels:    labels,
			Size:      size,
		})
	return m.save()
}

// RemoveDependency will remove the dependency with the given name from the
// untarred ACI stored at a.CurrentImagePath.
func (m *Manifest) RemoveDependency(imageName string) error {
	acid, err := types.NewACIdentifier(imageName)
	if err != nil {
		return err
	}

	err = removeDep(*acid, m.manifest)
	if err != nil {
		return err
	}
	return m.save()
}

func removeDep(imageName types.ACIdentifier, s *schema.ImageManifest) error {
	foundOne := false
	for i := len(s.Dependencies) - 1; i >= 0; i-- {
		if s.Dependencies[i].ImageName == imageName {
			foundOne = true
			s.Dependencies = append(
				s.Dependencies[:i],
				s.Dependencies[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
