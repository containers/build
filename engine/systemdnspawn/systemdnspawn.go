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

package systemdnspawn

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"

	"github.com/appc/spec/schema/types"
)

var pathlist = []string{"/usr/local/sbin", "/usr/local/bin", "/usr/sbin",
	"/usr/bin", "/sbin", "/bin"}

type Engine struct{}

func (e Engine) Run(command string, args []string, environment types.Environment, chroot, workingDir string) error {
	nspawncmd := []string{"systemd-nspawn", "-D", chroot}

	systemdVersion, err := getSystemdVersion()
	if err != nil {
		return err
	}

	if systemdVersion >= 209 {
		nspawncmd = append(nspawncmd, "--quiet", "--register=no")
	}
	if workingDir != "" {
		if systemdVersion < 229 {
			return fmt.Errorf("the working dir can only be set on systems with systemd-nspawn >= 229")
		}
		nspawncmd = append(nspawncmd, "--chdir", workingDir)
	}
	if systemdVersion >= 230 {
		machineIdFile := path.Join(chroot, "/etc/machine-id")
		_, err := os.Stat(machineIdFile)
		switch {
		case os.IsNotExist(err):
			err := os.MkdirAll(path.Dir(machineIdFile), 0755)
			if err != nil {
				return err
			}
			f, err := os.Create(machineIdFile)
			if err != nil {
				return err
			}
			f.Close()
			defer os.RemoveAll(path.Join(chroot, machineIdFile))
		case err != nil:
			return err
		}
	}

	for _, envVar := range environment {
		nspawncmd = append(nspawncmd, "--setenv", envVar.Name+"="+envVar.Value)
	}

	nspawncmd = append(nspawncmd, "--setenv", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")

	abscmd, err := findCmdInPath(pathlist, command, chroot)
	if err != nil {
		return err
	}

	finfo, err := os.Lstat(path.Join(chroot, abscmd))
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("binary %q doesn't exist", abscmd)
	case err != nil:
		return err
	}

	if finfo.Mode()&os.ModeSymlink != 0 && systemdVersion < 228 {
		fmt.Fprintf(os.Stderr, "Warning: %q is a symlink, which systemd-nspawn version %d might error on\n", abscmd, systemdVersion)
	}

	nspawncmd = append(nspawncmd, abscmd)
	nspawncmd = append(nspawncmd, args...)

	execCmd := exec.Command(nspawncmd[0], nspawncmd[1:]...)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = []string{"SYSTEMD_LOG_LEVEL=err"}

	err = execCmd.Run()
	if err == exec.ErrNotFound {
		return fmt.Errorf("systemd-nspawn is required but not found")
	}
	if err == nil {
		return nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		return fmt.Errorf("non-zero exit code: %d", code)
	}
	return err
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
