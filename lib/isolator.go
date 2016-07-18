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
	"encoding/json"
	"fmt"

	"github.com/appc/acbuild/util"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

func removeIsolatorFromMan(name types.ACIdentifier) func(*schema.ImageManifest) error {
	return func(s *schema.ImageManifest) error {
		if s.App == nil {
			return ErrNotFound
		}
		foundOne := false
		for i := len(s.App.Isolators) - 1; i >= 0; i-- {
			if s.App.Isolators[i].Name == name {
				foundOne = true
				s.App.Isolators = append(
					s.App.Isolators[:i],
					s.App.Isolators[i+1:]...)
			}
		}
		if !foundOne {
			return ErrNotFound
		}
		return nil
	}
}

func (a *ACBuild) AddIsolator(name string, value []byte) (err error) {
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
	rawMsg := json.RawMessage(value)

	fn := func(s *schema.ImageManifest) error {
		if s.App == nil {
			s.App = newManifestApp()
		}
		_, ok := types.ResourceIsolatorNames[*acid]
		if !ok {
			_, ok = types.LinuxIsolatorNames[*acid]
			if !ok {
				return fmt.Errorf("unknown isolator name: %s", name)
			}
		}
		i := &types.Isolator{
			Name:     *acid,
			ValueRaw: &rawMsg,
		}
		blob, err := json.Marshal(i)
		if err != nil {
			return err
		}
		err = i.UnmarshalJSON(blob)
		if err != nil {
			return err
		}
		removeIsolatorFromMan(*acid)(s)
		s.App.Isolators = append(s.App.Isolators, *i)
		return nil
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
}

func (a *ACBuild) RemoveIsolator(name string) (err error) {
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

	return util.ModifyManifest(removeIsolatorFromMan(*acid), a.CurrentACIPath)
}
