// Copyright 2016 The rkt Authors
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
	"fmt"
	"io"

	"github.com/appc/spec/schema/types"

	"github.com/containers/build/lib/appc"
)

// Manifest defines something that can manipulate manifests. The functions
// contained in it are the common subset of fields between appc and oci that can
// be altered by a user.
type Manifest interface {
	Print(w io.Writer, prettyPrint, printConfig bool) error // Print out this manifest to the given writer

	GetAnnotations() (map[string]string, error) // Used to generate build history

	AddAnnotation(name, value string) error
	AddEnv(name, value string) error
	AddLabel(name, value string) error
	AddMount(name, path string, readOnly bool) error
	AddPort(name, protocol string, port, count uint, socketActivated bool) error
	RemoveAnnotation(name string) error
	RemoveEnv(name string) error
	RemoveLabel(name string) error
	RemoveMount(name string) error
	RemovePort(name string) error
	Replace(manifestPath string) error
	SetExec(cmd []string) error
	SetGroup(group string) error
	SetUser(user string) error
	SetWorkingDir(dir string) error
	SetTag(tag string) error
}

func (a *ACBuild) Print(w io.Writer, prettyPrint, printConfig bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.Print(w, prettyPrint, printConfig)
}
func (a *ACBuild) GetAnnotations() (m map[string]string, err error) {
	if err = a.lock(); err != nil {
		return nil, err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.GetAnnotations()
}
func (a *ACBuild) AddAnnotation(name, value string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.AddAnnotation(name, value)
}
func (a *ACBuild) AddEnv(name, value string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.AddEnv(name, value)
}
func (a *ACBuild) AddLabel(name, value string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.AddLabel(name, value)
}
func (a *ACBuild) AddMount(name, path string, readOnly bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.AddMount(name, path, readOnly)
}
func (a *ACBuild) AddPort(name, protocol string, port, count uint, socketActivated bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.AddPort(name, protocol, port, count, socketActivated)
}
func (a *ACBuild) RemoveAnnotation(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.RemoveAnnotation(name)
}
func (a *ACBuild) RemoveEnv(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.RemoveEnv(name)
}
func (a *ACBuild) RemoveLabel(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.RemoveLabel(name)
}
func (a *ACBuild) RemoveMount(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.RemoveMount(name)
}
func (a *ACBuild) RemovePort(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.RemovePort(name)
}
func (a *ACBuild) Replace(manifestPath string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.Replace(manifestPath)
}
func (a *ACBuild) SetExec(cmd []string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.SetExec(cmd)
}
func (a *ACBuild) SetGroup(group string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.SetGroup(group)
}
func (a *ACBuild) SetUser(user string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.SetUser(user)
}
func (a *ACBuild) SetWorkingDir(dir string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.SetWorkingDir(dir)
}
func (a *ACBuild) SetTag(tag string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	return a.man.SetTag(tag)
}

func (a *ACBuild) AddDependency(imageName types.ACIdentifier, imageId *types.Hash, labels types.Labels, size uint) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.AddDependency(imageName, imageId, labels, size)
	}
	return fmt.Errorf("dependencies only supported in appc builds")
}
func (a *ACBuild) RemoveDependency(imageName string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.RemoveDependency(imageName)
	}
	return fmt.Errorf("dependencies only supported in appc builds")
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
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.AddIsolator(name, value)
	}
	return fmt.Errorf("dependencies only supported in appc builds")
}
func (a *ACBuild) RemoveIsolator(imageName string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.RemoveIsolator(imageName)
	}
	return fmt.Errorf("dependencies only supported in appc builds")
}
func (a *ACBuild) SetName(name string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.SetName(name)
	}
	return fmt.Errorf("setting names only supported in appc builds")
}
func (a *ACBuild) SetPostStop(exec []string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.SetPostStop(exec)
	}
	return fmt.Errorf("event handlers only supported in appc builds")
}
func (a *ACBuild) SetPreStart(exec []string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.SetPreStart(exec)
	}
	return fmt.Errorf("event handlers only supported in appc builds")
}
func (a *ACBuild) SetSuppGroups(groups []int) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()
	switch m := a.man.(type) {
	case *appc.Manifest:
		return m.SetSuppGroups(groups)
	}
	return fmt.Errorf("event handlers only supported in appc builds")
}
