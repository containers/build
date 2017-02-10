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

package chroot

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/containers/build/engine"
	"github.com/coreos/rkt/pkg/fileutil"
	"github.com/coreos/rkt/pkg/multicall"
	"github.com/coreos/rkt/pkg/user"
)

type Engine struct{}

func init() {
	multicall.Add("acbuild-chroot", cmdACBuildChroot.Execute)
}

func (e Engine) Run(command string, args []string, environment map[string]string, chroot, workingDir string) error {
	resolvConfFile := filepath.Join(chroot, "/etc/resolv.conf")
	_, err := os.Stat(resolvConfFile)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(filepath.Dir(resolvConfFile), 0755)
		if err != nil {
			return err
		}
		err = fileutil.CopyTree("/etc/resolv.conf", resolvConfFile, user.NewBlankUidRange())
		if err != nil {
			return err
		}
		defer os.RemoveAll(resolvConfFile)
	case err != nil:
		return err
	}
	var serializedArgs string
	for _, arg := range args {
		if serializedArgs != "" {
			serializedArgs += ","
		}
		serializedArgs += arg
	}
	var serializedEnv string
	for name, value := range environment {
		if serializedEnv != "" {
			serializedEnv += ","
		}
		serializedEnv += name + "=" + value
	}
	path := "PATH="
	for _, p := range engine.Pathlist {
		if path != "PATH=" {
			path += ":"
		}
		path += p
	}
	chrootArgs := []string{
		"--cmd", command,
		"--chroot", chroot,
		"--working-dir", workingDir,
	}
	if len(serializedArgs) > 0 {
		chrootArgs = append(chrootArgs, "--args", serializedArgs)
	}
	if len(serializedEnv) > 0 {
		chrootArgs = append(chrootArgs, "--env", serializedEnv)
	}
	cmd := exec.Command("acbuild-chroot", chrootArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = []string{path}
	return cmd.Run()
}
