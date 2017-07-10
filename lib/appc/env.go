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
)

// AddEnv will add an environment variable of name and value to the current
// manifest
func (m *Manifest) AddEnv(name, value string) error {
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	m.manifest.App.Environment.Set(name, value)
	return m.save()
}

// Remove Env will remove an environment variable of name from the current
// manifest
func (m *Manifest) RemoveEnv(name string) error {
	err := removeFromEnv(name, m.manifest)
	if err != nil {
		return err
	}
	return m.save()
}

func removeFromEnv(name string, m *schema.ImageManifest) error {
	if m.App == nil {
		return ErrNotFound
	}
	foundOne := false
	for i := len(m.App.Environment) - 1; i >= 0; i-- {
		if m.App.Environment[i].Name == name {
			foundOne = true
			m.App.Environment = append(
				m.App.Environment[:i],
				m.App.Environment[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
