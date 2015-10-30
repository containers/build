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

const (
	preStartName = "pre-start"
	postStopName = "post-stop"
)

// SetPreStart sets the pre-start event handler in the expanded ACI stored at
// a.CurrentACIPath
func (a *ACBuild) SetPreStart(exec []string) error {
	return a.setEventHandler(preStartName, exec)
}

// SetPostStop sets the post-stop event handler in the expanded ACI stored at
// a.CurrentACIPath
func (a *ACBuild) SetPostStop(exec []string) error {
	return a.setEventHandler(postStopName, exec)
}

func (a *ACBuild) setEventHandler(name string, exec []string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	fn := func(s *schema.ImageManifest) {
		removeEventHandler(name, s)
		if s.App == nil {
			s.App = &types.App{}
		}
		s.App.EventHandlers = append(s.App.EventHandlers,
			types.EventHandler{
				Name: name,
				Exec: exec,
			})
	}
	return util.ModifyManifest(fn, a.CurrentACIPath)
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
