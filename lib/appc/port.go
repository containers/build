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
	"fmt"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

// AddPort will add a port with the given name, protocol, port, and count to
// the untarred ACI stored at a.CurrentImagePath. If the port already exists its
// value will be updated to the new value. socketActivated signifies whether or
// not the application will be socket activated via this port.
func (m *Manifest) AddPort(name, protocol string, port, count uint, socketActivated bool) error {
	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	removePort(name, m.manifest)
	m.manifest.App.Ports = append(m.manifest.App.Ports,
		types.Port{
			Name:            *acn,
			Protocol:        protocol,
			Port:            port,
			Count:           count,
			SocketActivated: socketActivated,
		})

	return m.save()
}

// RemovePort will remove the port with the given name from the untarred ACI
// stored at a.CurrentImagePath.
func (m *Manifest) RemovePort(port string) error {
	err := removePort(port, m.manifest)
	if err != nil {
		return err
	}

	return m.save()
}

func removePort(port string, s *schema.ImageManifest) error {
	if s.App == nil {
		return ErrNotFound
	}
	acn, err := types.NewACName(port)
	if err != nil {
		return err
	}
	foundOne := false
	for i := len(s.App.Ports) - 1; i >= 0; i-- {
		if s.App.Ports[i].Name == *acn || fmt.Sprintf("%d", s.App.Ports[i].Port) == port {
			foundOne = true
			s.App.Ports = append(
				s.App.Ports[:i],
				s.App.Ports[i+1:]...)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	return nil
}
