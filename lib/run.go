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

	"github.com/containers/build/engine"
	"github.com/containers/build/lib/oci"
	"github.com/containers/build/registry"
	"github.com/containers/build/util"
)

// Run will execute the given command in the ACI being built. a.CurrentImagePath
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

	var depPaths []string
	switch a.Mode {
	case BuildModeOCI:
		depPaths, err = a.generateOverlayPathsOCI()
	case BuildModeAppC:
		depPaths, err = a.generateOverlayPathsAppC(insecure)
	default:
		return fmt.Errorf("unknown build mode: %s", a.Mode)
	}
	if err != nil {
		return err
	}

	if len(depPaths) != 1 {
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

	var chrootDir string
	if len(depPaths) == 1 {
		chrootDir = depPaths[0]
	} else {
		lowerLayers := depPaths[0 : len(depPaths)-1]
		upperLayer := depPaths[len(depPaths)-1]
		options := "lowerdir=" + strings.Join(lowerLayers, ":") +
			",upperdir=" + upperLayer +
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

	var env map[string]string
	switch a.Mode {
	case BuildModeOCI:
		env, err = a.getEnvVarsOCI()
	case BuildModeAppC:
		env, err = a.getEnvVarsAppC()
	default:
		return fmt.Errorf("unknown build mode: %s", a.Mode)
	}
	if err != nil {
		return err
	}

	err = a.mirrorLocalZoneInfo()
	if err != nil {
		return err
	}

	err = runEngine.Run(cmd[0], cmd[1:], env, chrootDir, workingDir)
	if err != nil {
		return err
	}

	if a.Mode == BuildModeOCI {
		err = a.rehashAndStoreOCIBlob(depPaths[len(depPaths)-1], false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *ACBuild) generateOverlayPathsAppC(insecure bool) ([]string, error) {
	deps, err := a.renderACI(insecure, a.Debug)
	if err != nil {
		return nil, err
	}
	for i, dep := range deps {
		deps[i] = path.Join(a.DepStoreExpandedPath, dep, aci.RootfsDir)
	}

	deps = append(deps, path.Join(a.CurrentImagePath, aci.RootfsDir))

	return deps, nil
}

func (a *ACBuild) generateOverlayPathsOCI() ([]string, error) {
	var layerIDs []string
	switch ociMan := a.man.(type) {
	case *oci.Image:
		layerIDs = ociMan.GetLayerHashes()
	default:
		return nil, fmt.Errorf("internal error: mismatched manifest type and build mode???")
	}

	var layerPaths []string
	if len(layerIDs) == 0 {
		layerPaths = []string{path.Join(a.OCIExpandedBlobsPath, "sha256", "new-layer")}
		err := os.MkdirAll(layerPaths[0], 0755)
		if err != nil {
			return nil, err
		}
	} else {
		err := util.OCIExtractLayers(layerIDs, a.CurrentImagePath, a.OCIExpandedBlobsPath)
		if err != nil {
			return nil, err
		}
		for _, layerID := range layerIDs {
			algo, hash, err := util.SplitOCILayerID(layerID)
			if err != nil {
				return nil, err
			}
			layerPaths = append(layerPaths, path.Join(a.OCIExpandedBlobsPath, algo, hash))
		}
	}
	return layerPaths, nil
}

func (a *ACBuild) getEnvVarsAppC() (map[string]string, error) {
	man, err := util.GetManifest(a.CurrentImagePath)
	if err != nil {
		return nil, err
	}

	var env types.Environment
	if man.App != nil {
		env = man.App.Environment
	} else {
		env = types.Environment{}
	}

	envMap := make(map[string]string)
	for _, v := range env {
		envMap[v.Name] = v.Value
	}

	return envMap, nil
}

func (a *ACBuild) getEnvVarsOCI() (map[string]string, error) {
	switch ociMan := a.man.(type) {
	case *oci.Image:
		env := ociMan.GetConfig().Config.Env
		ret := make(map[string]string)
		for _, v := range env {
			tokens := strings.SplitN(v, "=", 2)
			if len(tokens) < 2 {
				return nil, fmt.Errorf("incorrectly formatted environment variable: %q", v)
			}
			ret[tokens[0]] = tokens[1]
		}
		return ret, nil
	default:
		return nil, fmt.Errorf("internal error: mismatched manifest type and build mode???")
	}
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

	man, err := util.GetManifest(a.CurrentImagePath)
	if err != nil {
		return nil, err
	}

	if len(man.Dependencies) == 0 {
		return nil, nil
	}

	var deplist []string
	for _, dep := range man.Dependencies {
		err := reg.FetchAndRender(dep.ImageName, dep.Labels, dep.Size)
		switch err {
		case nil:
		case registry.ErrNotFound:
			l, _ := dep.Labels.Get("version")
			return nil, fmt.Errorf("dependency %q doesn't appear to exist: %v", string(dep.ImageName)+":"+l, err)
		default:
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

	destp := filepath.Join(a.CurrentImagePath, aci.RootfsDir, zif)

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
