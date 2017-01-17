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
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/appc/spec/schema/types"
)

const defaultWorkPath = ".acbuild"

var (
	// ErrNotFound is returned when acbuild is asked to remove an element from a
	// list and the element is not present in the list
	ErrNotFound = fmt.Errorf("element to be removed does not exist in this ACI")

	errNoBuildInProgress = fmt.Errorf("no build in progress in this working dir - try \"acbuild begin\"")
)

// newManifestApp will generate a valid minimal types.App for use in a
// schema.ImageManifest. This is necessary as placing a completely empty
// types.App into a manifest will result in an invalid manifest.
func newManifestApp() *types.App {
	return &types.App{
		User:  "0",
		Group: "0",
	}
}

// ACBuild contains all the information for a current build. Once an ACBuild
// has been created, the functions available on it will perform different
// actions in the build, like updating a dependency or writing a finished ACI.
type ACBuild struct {
	ContextPath          string
	LockPath             string
	CurrentACIPath       string
	DepStoreTarPath      string
	DepStoreExpandedPath string
	OverlayTargetPath    string
	OverlayWorkPath      string
	Debug                bool

	lockFile *os.File
}

// NewACBuild returns a new ACBuild struct with sane defaults for all of the
// different paths
func NewACBuild(cwd string, debug bool) *ACBuild {
	return &ACBuild{
		ContextPath:          path.Join(cwd, defaultWorkPath),
		LockPath:             path.Join(cwd, defaultWorkPath, "lock"),
		CurrentACIPath:       path.Join(cwd, defaultWorkPath, "currentaci"),
		DepStoreTarPath:      path.Join(cwd, defaultWorkPath, "depstore-tar"),
		DepStoreExpandedPath: path.Join(cwd, defaultWorkPath, "depstore-expanded"),
		OverlayTargetPath:    path.Join(cwd, defaultWorkPath, "target"),
		OverlayWorkPath:      path.Join(cwd, defaultWorkPath, "work"),
		Debug:                debug,
	}
}

func (a *ACBuild) lock() error {
	_, err := os.Stat(a.ContextPath)
	switch {
	case os.IsNotExist(err):
		return errNoBuildInProgress
	case err != nil:
		return err
	}

	if a.lockFile != nil {
		return fmt.Errorf("lock already held by this ACBuild")
	}

	a.lockFile, err = os.OpenFile(a.LockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		switch err1 := err.(type) {
		case *os.PathError:
			switch err2 := err1.Err.(type) {
			case syscall.Errno:
				if err2 == syscall.EACCES {
					err = fmt.Errorf("permission denied: please run this as a user with appropriate privileges\n")
				}
			}
		}
		return err
	}

	err = syscall.Flock(int(a.lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		if err == syscall.EWOULDBLOCK {
			return fmt.Errorf("lock already held - is another acbuild running in this working dir?")
		}
		return err
	}

	return nil
}

func (a *ACBuild) unlock() error {
	if a.lockFile == nil {
		return fmt.Errorf("lock isn't held by this ACBuild")
	}

	err := syscall.Flock(int(a.lockFile.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}

	err = a.lockFile.Close()
	if err != nil {
		return err
	}
	a.lockFile = nil

	err = os.Remove(a.LockPath)
	if err != nil {
		return err
	}

	return nil
}
