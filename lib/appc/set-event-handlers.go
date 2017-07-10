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

const (
	preStartName = "pre-start"
	postStopName = "post-stop"
)

// SetPreStart sets the pre-start event handler in the expanded ACI stored at
// a.CurrentImagePath
func (m *Manifest) SetPreStart(exec []string) error {
	return m.setEventHandler(preStartName, exec)
}

// SetPostStop sets the post-stop event handler in the expanded ACI stored at
// a.CurrentImagePath
func (m *Manifest) SetPostStop(exec []string) error {
	return m.setEventHandler(postStopName, exec)
}

func (m *Manifest) setEventHandler(name string, exec []string) (err error) {
	removeEventHandler(name, m.manifest)
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	m.manifest.App.EventHandlers = append(m.manifest.App.EventHandlers,
		types.EventHandler{
			Name: name,
			Exec: exec,
		})
	return m.save()
}

func removeEventHandler(name string, s *schema.ImageManifest) {
	if s.App == nil {
		return
	}
	for i := len(s.App.EventHandlers) - 1; i >= 0; i-- {
		if s.App.EventHandlers[i].Name == name {
			s.App.EventHandlers = append(
				s.App.EventHandlers[:i],
				s.App.EventHandlers[i+1:]...)
		}
	}
}
