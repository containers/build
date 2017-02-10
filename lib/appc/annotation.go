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

// AddAnnotation adds an annotation of name and value to the current manifest
func (m *Manifest) AddAnnotation(name, value string) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}
	m.manifest.Annotations.Set(*acid, value)
	return m.save()
}

// RemoveAnnotation removes an annotation of name from the current manifest
func (m *Manifest) RemoveAnnotation(name string) error {
	acid, err := types.NewACIdentifier(name)
	if err != nil {
		return err
	}

	err = removeAnnotation(*acid, m.manifest)
	if err != nil {
		return err
	}
	return m.save()
}

func removeAnnotation(name types.ACIdentifier, s *schema.ImageManifest) error {
	foundOne := false
	for i := len(s.Annotations) - 1; i >= 0; i-- {
		if s.Annotations[i].Name == name {
			foundOne = true
			s.Annotations = append(
				s.Annotations[:i],
				s.Annotations[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}

func (m *Manifest) GetAnnotations() (map[string]string, error) {
	ret := make(map[string]string)
	for _, v := range m.manifest.Annotations {
		ret[string(v.Name)] = v.Value
	}
	return ret, nil
}
