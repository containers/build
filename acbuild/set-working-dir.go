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

package main

import (
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	cmdSetWorkingDir = &cobra.Command{
		Use:     "set-working-directory DIR",
		Short:   "Set the working directory",
		Long:    "Set the working directory the app will run in inside the container",
		Example: "acbuild set-working-directory /root",
		Aliases: []string{"set-wd"},
		Run:     runWrapper(runSetWorkingDir),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdSetWorkingDir)
}

func runSetWorkingDir(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("set-working-dir: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Setting working directory to %s", args[0])
	}

	err := newACBuild().SetWorkingDir(args[0])

	if err != nil {
		stderr("set-working-dir: %v", err)
		return 1
	}

	return 0
}
