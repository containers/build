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
)

func removeFromEnv(name string) func(*schema.ImageManifest) error {
	return func(s *schema.ImageManifest) error {
		if s.App == nil {
			return ErrNotFound
		}
		foundOne := false
		for i := len(s.App.Environment) - 1; i >= 0; i-- {
			if s.App.Environment[i].Name == name {
				foundOne = true
				s.App.Environment = append(
					s.App.Environment[:i],
					s.App.Environment[i+1:]...)
			}
		}
		if !foundOne {
			return ErrNotFound
		}
		return nil
	}
}

// AddEnv will add an environment variable with the given name and value to the
// untarred ACI stored at a.CurrentACIPath. If the environment variable already
// exists its value will be updated to the new value.
func (a *ACBuild) AddEnv(name, value string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	fn := func(s *schema.ImageManifest) error {
		if s.App == nil {
			s.App = newManifestApp()
		}
		s.App.Environment.Set(name, value)
		return nil
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
}

// RemoveEnv will remove the environment variable with the given name from the
// untarred ACI stored at a.CurrentACIPath.
func (a *ACBuild) RemoveEnv(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	return util.ModifyManifest(removeFromEnv(name), a.CurrentACIPath)
}
