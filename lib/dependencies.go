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

package lib

import (
	"github.com/appc/acbuild/util"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

func removeDep(imageName types.ACIdentifier) func(*schema.ImageManifest) error {
	return func(s *schema.ImageManifest) error {
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
}

// AddDependency will add a dependency with the given name, id, labels, and size
// to the untarred ACI stored at a.CurrentACIPath. If the dependency already
// exists its fields will be updated to the new values.
func (a *ACBuild) AddDependency(imageName types.ACIdentifier, imageId *types.Hash, labels types.Labels, size uint) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	fn := func(s *schema.ImageManifest) error {
		removeDep(imageName)(s)
		s.Dependencies = append(s.Dependencies,
			types.Dependency{
				ImageName: imageName,
				ImageID:   imageId,
				Labels:    labels,
				Size:      size,
			})
		return nil
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
}

// RemoveDependency will remove the dependency with the given name from the
// untarred ACI stored at a.CurrentACIPath.
func (a *ACBuild) RemoveDependency(imageName string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	acid, err := types.NewACIdentifier(imageName)
	if err != nil {
		return err
	}

	return util.ModifyManifest(removeDep(*acid), a.CurrentACIPath)
}
