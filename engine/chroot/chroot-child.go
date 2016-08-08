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
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

func init() {
	cmdACBuildChroot.PersistentFlags().StringVar(&flagCmd, "cmd", "", "The command to run")
	cmdACBuildChroot.PersistentFlags().StringSliceVar(&flagArgs, "args", nil, "arguments for the command")
	cmdACBuildChroot.PersistentFlags().StringSliceVar(&flagEnv, "env", nil, "environment for the command")
	cmdACBuildChroot.PersistentFlags().StringVar(&flagChroot, "chroot", "", "dir to chroot into")
	cmdACBuildChroot.PersistentFlags().StringVar(&flagWorkingDir, "working-dir", "", "working directory for the command")
}

var (
	flagCmd          string
	flagArgs         []string
	flagEnv          []string
	flagChroot       string
	flagWorkingDir   string
	cmdACBuildChroot = &cobra.Command{
		Use: "",
		Run: runChroot,
	}
)

func stderr(format string, a ...interface{}) {
	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
}

func errAndExit(format string, a ...interface{}) {
	stderr(format, a...)
	os.Exit(1)
}

func runChroot(cmd *cobra.Command, args []string) {
	runtime.LockOSThread()
	err := syscall.Chroot(flagChroot)
	if err != nil {
		errAndExit("couldn't chroot: %v", err)
	}
	err = os.Chdir("/")
	if err != nil {
		errAndExit("couldn't cd: %v", err)
	}

	if flagWorkingDir != "" {
		err = os.Chdir(flagWorkingDir)
		if err != nil {
			errAndExit("couldn't cd: %v", err)
		}
	}

	execCmd := exec.Command(flagCmd, flagArgs...)
	execCmd.Env = flagEnv
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	err = execCmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
			os.Exit(code)
		}
		errAndExit("%v", err)
	}
}
