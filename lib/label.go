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

func removeLabelFromMan(name types.ACIdentifier) func(*schema.ImageManifest) error {
	return func(s *schema.ImageManifest) error {
		foundOne := false
		for i := len(s.Labels) - 1; i >= 0; i-- {
			if s.Labels[i].Name == name {
				foundOne = true
				s.Labels = append(
					s.Labels[:i],
					s.Labels[i+1:]...)
			}
		}
		if !foundOne {
			return ErrNotFound
		}
		return nil
	}
}

// AddLabel will add a label with the given name and value to the untarred ACI
// stored at a.CurrentACIPath. If the label already exists its value will be updated to
// the new value.
func (a *ACBuild) AddLabel(name, value string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	fn := func(s *schema.ImageManifest) error {
		removeLabelFromMan(*acid)(s)
		s.Labels = append(s.Labels,
			types.Label{
				Name:  *acid,
				Value: value,
			})
		return nil
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
}

// RemoveLabel will remove the label with the given name from the untarred ACI
// stored at a.CurrentACIPath
func (a *ACBuild) RemoveLabel(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	return util.ModifyManifest(removeLabelFromMan(*acid), a.CurrentACIPath)
}
