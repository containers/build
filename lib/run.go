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
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/schema/types"

	"github.com/appc/acbuild/engine"
	"github.com/appc/acbuild/registry"
	"github.com/appc/acbuild/util"
)

// Run will execute the given command in the ACI being built. a.CurrentACIPath
// is where the untarred ACI is stored, a.DepStoreTarPath is the directory to
// download dependencies into, a.DepStoreExpandedPath is where the dependencies
// are expanded into, and a.OverlayWorkPath is the work directory used by
// overlayfs.
//
// Arguments:
//
// - cmd:        The command to run and its arguments.
//
// - workingDir: If specified, the current directory inside the container is
// changed to its value before running the given command.
//
// - runEngine:  The engine used to perform the execution of the command.
func (a *ACBuild) Run(cmd []string, workingDir string, insecure bool, runEngine engine.Engine) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	if os.Geteuid() != 0 {
		return fmt.Errorf("the run subcommand must be run as root")
	}

	if len(cmd) == 0 {
		return fmt.Errorf("command to run not set")
	}

	err = util.MaybeUnmount(a.OverlayTargetPath)
	if err != nil {
		return err
	}

	err = util.RmAndMkdir(a.OverlayTargetPath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(a.OverlayTargetPath)
	err = util.RmAndMkdir(a.OverlayWorkPath)
	if err != nil {
		return err
	}
	defer os.RemoveAll(a.OverlayWorkPath)
	err = os.MkdirAll(a.DepStoreExpandedPath, 0755)
	if err != nil {
		return err
	}
	err = os.MkdirAll(a.DepStoreTarPath, 0755)
	if err != nil {
		return err
	}

	man, err := util.GetManifest(a.CurrentACIPath)
	if err != nil {
		return err
	}

	if len(man.Dependencies) != 0 {
		if !supportsOverlay() {
			err := exec.Command("modprobe", "overlay").Run()
			if err != nil {
				if _, ok := err.(*exec.ExitError); ok {
					return fmt.Errorf("overlayfs is not supported on your system")
				}
				return err
			}
			if !supportsOverlay() {
				return fmt.Errorf(
					"overlayfs support required for using run with dependencies")
			}
		}
	}

	deps, err := a.renderACI(insecure, a.Debug)
	if err != nil {
		return err
	}

	var chrootDir string
	if deps == nil {
		chrootDir = path.Join(a.CurrentACIPath, aci.RootfsDir)
	} else {
		for i, dep := range deps {
			deps[i] = path.Join(a.DepStoreExpandedPath, dep, aci.RootfsDir)
		}
		options := "lowerdir=" + strings.Join(deps, ":") +
			",upperdir=" + path.Join(a.CurrentACIPath, aci.RootfsDir) +
			",workdir=" + a.OverlayWorkPath
		err := syscall.Mount("overlay", a.OverlayTargetPath, "overlay", 0, options)
		if err != nil {
			return err
		}

		defer func() {
			err1 := syscall.Unmount(a.OverlayTargetPath, 0)
			if err == nil {
				err = err1
			}
		}()

		chrootDir = a.OverlayTargetPath
	}

	var env types.Environment
	if man.App != nil {
		env = man.App.Environment
	} else {
		env = types.Environment{}
	}

	err = a.mirrorLocalZoneInfo()
	if err != nil {
		return err
	}

	err = runEngine.Run(cmd[0], cmd[1:], env, chrootDir, workingDir)
	if err != nil {
		return err
	}

	return nil
}

// stolen from github.com/coreos/rkt/common/common.go
// supportsOverlay returns whether the system supports overlay filesystem
func supportsOverlay() bool {
	f, err := os.Open("/proc/filesystems")
	if err != nil {
		fmt.Println("error opening /proc/filesystems")
		return false
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		if s.Text() == "nodev\toverlay" {
			return true
		}
	}
	return false
}

func (a *ACBuild) renderACI(insecure, debug bool) ([]string, error) {
	reg := registry.Registry{
		DepStoreTarPath:      a.DepStoreTarPath,
		DepStoreExpandedPath: a.DepStoreExpandedPath,
		Insecure:             insecure,
		Debug:                debug,
	}

	man, err := util.GetManifest(a.CurrentACIPath)
	if err != nil {
		return nil, err
	}

	if len(man.Dependencies) == 0 {
		return nil, nil
	}

	var deplist []string
	for _, dep := range man.Dependencies {
		err := reg.FetchAndRender(dep.ImageName, dep.Labels, dep.Size)
		if err != nil {
			return nil, err
		}

		depkey, err := reg.GetACI(dep.ImageName, dep.Labels)
		if err != nil {
			return nil, err
		}

		subdeplist, err := genDeplist(path.Join(a.DepStoreExpandedPath, depkey), reg)
		if err != nil {
			return nil, err
		}
		deplist = append(deplist, subdeplist...)
	}

	return deplist, nil
}

func genDeplist(acipath string, reg registry.Registry) ([]string, error) {
	man, err := util.GetManifest(acipath)
	if err != nil {
		return nil, err
	}
	key, err := reg.GetACI(man.Name, man.Labels)
	if err != nil {
		fmt.Printf("Name: %s", man.Name)
		return nil, err
	}

	var deps []string
	for _, dep := range man.Dependencies {
		depkey, err := reg.GetACI(dep.ImageName, dep.Labels)
		if err != nil {
			return nil, err
		}

		subdeps, err := genDeplist(path.Join(reg.DepStoreExpandedPath, depkey), reg)
		if err != nil {
			return nil, err
		}
		deps = append(deps, subdeps...)
	}

	deps = append(deps, key)
	return deps, nil
}

func (a *ACBuild) mirrorLocalZoneInfo() error {
	zif, err := filepath.EvalSymlinks("/etc/localtime")
	if err != nil {
		return err
	}

	src, err := os.Open(zif)
	if err != nil {
		return err
	}
	defer src.Close()

	destp := filepath.Join(a.CurrentACIPath, aci.RootfsDir, zif)

	if err = os.MkdirAll(filepath.Dir(destp), 0755); err != nil {
		return err
	}

	dest, err := os.OpenFile(destp, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}
