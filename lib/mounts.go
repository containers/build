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

func removeMount(name types.ACName) func(*schema.ImageManifest) {
	return func(s *schema.ImageManifest) {
		if s.App == nil {
			return
		}
		for i := len(s.App.MountPoints) - 1; i >= 0; i-- {
			if s.App.MountPoints[i].Name == name {
				s.App.MountPoints = append(
					s.App.MountPoints[:i],
					s.App.MountPoints[i+1:]...)
			}
		}
	}
}

// AddMount will add a mount point with the given name and path to the untarred
// ACI stored at acipath. If the mount point already exists its value will be
// updated to the new value. readOnly signifies whether or not the mount point
// should be read only.
func AddMount(acipath, name, path string, readOnly bool) error {
	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}

	fn := func(s *schema.ImageManifest) {
		removeMount(*acn)(s)
		if s.App == nil {
			s.App = &types.App{}
		}
		s.App.MountPoints = append(s.App.MountPoints,
			types.MountPoint{
				Name:     *acn,
				Path:     path,
				ReadOnly: readOnly,
			})
	}
	return util.ModifyManifest(fn, acipath)
}

// RemoveMount will remove the mount point with the given name from the
// untarred ACI stored at acipath
func RemoveMount(acipath, name string) error {
	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}

	return util.ModifyManifest(removeMount(*acn), acipath)
}
