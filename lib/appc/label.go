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

package appc

import (
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

func (m *Manifest) SetTag(value string) error {
	return m.AddLabel("version", value)
}

// AddLabel will add a label with the given name and value to the untarred ACI
// stored at a.CurrentImagePath. If the label already exists its value will be updated to
// the new value.
func (m *Manifest) AddLabel(name, value string) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	removeLabelFromMan(*acid, m.manifest)
	m.manifest.Labels = append(m.manifest.Labels,
		types.Label{
			Name:  *acid,
			Value: value,
		})

	return m.save()
}

// RemoveLabel will remove the label with the given name from the untarred ACI
// stored at a.CurrentImagePath
func (m *Manifest) RemoveLabel(name string) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	err = removeLabelFromMan(*acid, m.manifest)
	if err != nil {
		return err
	}

	return m.save()
}

func removeLabelFromMan(name types.ACIdentifier, m *schema.ImageManifest) error {
	foundOne := false
	for i := len(m.Labels) - 1; i >= 0; i-- {
		if m.Labels[i].Name == name {
			foundOne = true
			m.Labels = append(
				m.Labels[:i],
				m.Labels[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
