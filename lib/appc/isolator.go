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
	"encoding/json"
	"fmt"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

// AddIsolator adds an isolator of name and value to the current manifest
func (m *Manifest) AddIsolator(name string, value []byte) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}
	rawMsg := json.RawMessage(value)

	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
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
	removeIsolatorFromMan(*acid, m.manifest)
	m.manifest.App.Isolators = append(m.manifest.App.Isolators, *i)
	return m.save()
}

// AddIsolator removes an isolator of name from the current manifest
func (m *Manifest) RemoveIsolator(name string) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	err = removeIsolatorFromMan(*acid, m.manifest)
	if err != nil {
		return err
	}
	return m.save()
}

func removeIsolatorFromMan(name types.ACIdentifier, m *schema.ImageManifest) error {
	if m.App == nil {
		return ErrNotFound
	}
	foundOne := false
	for i := len(m.App.Isolators) - 1; i >= 0; i-- {
		if m.App.Isolators[i].Name == name {
			foundOne = true
			m.App.Isolators = append(
				m.App.Isolators[:i],
				m.App.Isolators[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
