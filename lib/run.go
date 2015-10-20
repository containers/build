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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"

	"github.com/appc/acbuild/registry"
	"github.com/appc/acbuild/util"
)

var pathlist = []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin",
	"/usr/bin", "/sbin", "/bin"}

// Run will execute the given command in the ACI being built. a.CurrentACIPath
// is where the untarred ACI is stored, a.DepStoreTarPath is the directory to
// download dependencies into, a.DepStoreExpandedPath is where the dependencies
// are expanded into, a.OverlayWorkPath is the work directory used by
// overlayfs, and insecure signifies whether downloaded images should be
// fetched over http or https.
func (a *ACBuild) Run(cmd []string, insecure bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

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
			err := util.Exec("modprobe", "overlay")
			if err != nil {
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

	var nspawnpath string
	if deps == nil {
		nspawnpath = path.Join(a.CurrentACIPath, aci.RootfsDir)
	} else {
		for i, dep := range deps {
			deps[i] = path.Join(a.DepStoreExpandedPath, dep, aci.RootfsDir)
		}
		options := "-olowerdir=" + strings.Join(deps, ":") +
			",upperdir=" + path.Join(a.CurrentACIPath, aci.RootfsDir) +
			",workdir=" + a.OverlayWorkPath
		err := util.Exec("mount", "-t", "overlay",
			"overlay", options, a.OverlayTargetPath)
		if err != nil {
			return err
		}

		umount := exec.Command("umount", a.OverlayTargetPath)
		umount.Stdout = os.Stdout
		umount.Stderr = os.Stderr
		defer umount.Run()

		nspawnpath = a.OverlayTargetPath
	}
	nspawncmd := []string{"systemd-nspawn", "-q", "-D", nspawnpath}

	if man.App != nil {
		for _, evar := range man.App.Environment {
			nspawncmd = append(nspawncmd, "--setenv", evar.Name+"="+evar.Value)
		}
	}

	if len(cmd) == 0 {
		return fmt.Errorf("command to run not set")
	}
	abscmd, err := findCmdInPath(pathlist, cmd[0], nspawnpath)
	if err != nil {
		return err
	}
	nspawncmd = append(nspawncmd, abscmd)
	nspawncmd = append(nspawncmd, cmd[1:]...)
	//fmt.Printf("%v\n", nspawncmd)

	err = util.Exec(nspawncmd[0], nspawncmd[1:]...)
	if err != nil {
		if err == exec.ErrNotFound {
			return fmt.Errorf("systemd-nspawn is required but not found")
		}
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

func findCmdInPath(pathlist []string, cmd, prefix string) (string, error) {
	if len(cmd) != 0 && cmd[0] == '/' {
		return cmd, nil
	}

	for _, p := range pathlist {
		ex, err := util.Exists(path.Join(prefix, p, cmd))
		if err != nil {
			return "", err
		}
		if ex {
			return path.Join(p, cmd), nil
		}
	}
	return "", fmt.Errorf("%s not found in any of: %v", cmd, pathlist)
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
