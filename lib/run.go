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

	if os.Geteuid() != 0 {
		return fmt.Errorf("the run subcommand must be run as root")
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

	var nspawnpath string
	if deps == nil {
		nspawnpath = path.Join(a.CurrentACIPath, aci.RootfsDir)
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

		nspawnpath = a.OverlayTargetPath
	}
	nspawncmd := []string{"systemd-nspawn", "-D", nspawnpath}

	systemdVersion, err := getSystemdVersion()
	if err != nil {
		return err
	}
	if systemdVersion >= 209 {
		nspawncmd = append(nspawncmd, "--quiet", "--register=no")
	}

	if man.App != nil {
		for _, evar := range man.App.Environment {
			nspawncmd = append(nspawncmd, "--setenv", evar.Name+"="+evar.Value)
		}
	}
	nspawncmd = append(nspawncmd, "--setenv", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	err = a.mirrorLocalZoneInfo()
	if err != nil {
		return err
	}

	if len(cmd) == 0 {
		return fmt.Errorf("command to run not set")
	}
	abscmd, err := findCmdInPath(pathlist, cmd[0], nspawnpath)
	if err != nil {
		return err
	}

	finfo, err := os.Lstat(path.Join(nspawnpath, abscmd))
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("the binary %q doesn't exist", abscmd)
	case err != nil:
		return err
	case finfo.Mode()&os.ModeSymlink != 0 && systemdVersion < 228:
		fmt.Fprintf(os.Stderr, "Warning: %q is a symlink, which systemd-nspawn version %d might error on\n", abscmd, systemdVersion)
	}

	nspawncmd = append(nspawncmd, abscmd)
	nspawncmd = append(nspawncmd, cmd[1:]...)

	execCmd := exec.Command(nspawncmd[0], nspawncmd[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = []string{"SYSTEMD_LOG_LEVEL=err"}

	err = execCmd.Run()
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
	if path.IsAbs(cmd) {
		return cmd, nil
	}

	for _, p := range pathlist {
		_, err := os.Lstat(path.Join(prefix, p, cmd))
		switch {
		case os.IsNotExist(err):
			continue
		case err != nil:
			return "", err
		}
		return path.Join(p, cmd), nil
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

func getSystemdVersion() (int, error) {
	_, err := exec.LookPath("systemctl")
	if err == exec.ErrNotFound {
		return 0, fmt.Errorf("system does not have systemd")
	}

	blob, err := exec.Command("systemctl", "--version").Output()
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(string(blob), "\n") {
		if strings.HasPrefix(line, "systemd ") {
			var version int
			_, err := fmt.Sscanf(line, "systemd %d", &version)
			if err != nil {
				return 0, err
			}
			return version, nil
		}
	}
	return 0, fmt.Errorf("error parsing output from `systemctl --version`")
}
