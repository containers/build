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

func removePort(name types.ACName) func(*schema.ImageManifest) error {
	return func(s *schema.ImageManifest) error {
		if s.App == nil {
			return nil
		}
		foundOne := false
		for i := len(s.App.Ports) - 1; i >= 0; i-- {
			if s.App.Ports[i].Name == name {
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
}

// AddPort will add a port with the given name, protocol, port, and count to
// the untarred ACI stored at a.CurrentACIPath. If the port already exists its
// value will be updated to the new value. socketActivated signifies whether or
// not the application will be socket activated via this port.
func (a *ACBuild) AddPort(name, protocol string, port, count uint, socketActivated bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}

	fn := func(s *schema.ImageManifest) error {
		removePort(*acn)(s)
		if s.App == nil {
			s.App = &types.App{}
		}
		s.App.Ports = append(s.App.Ports,
			types.Port{
				Name:            *acn,
				Protocol:        protocol,
				Port:            port,
				Count:           count,
				SocketActivated: socketActivated,
			})
		return nil
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
}

// RemovePort will remove the port with the given name from the untarred ACI
// stored at a.CurrentACIPath.
func (a *ACBuild) RemovePort(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	acn, err := types.NewACName(name)
	if err != nil {
		return err
	}

	return util.ModifyManifest(removePort(*acn), a.CurrentACIPath)
}
