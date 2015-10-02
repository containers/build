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

	"github.com/appc/acbuild/lib"
)

var (
	cmdEnv = &cobra.Command{
		Use:     "environment [command]",
		Aliases: []string{"env"},
		Short:   "Manage environment variables",
	}
	cmdAddEnv = &cobra.Command{
		Use:     "add NAME VALUE",
		Short:   "Add an environment variable",
		Long:    "Updates the ACI to contain an environment variable with the given name and value. If the variable already exists, its value will be changed.",
		Example: "acbuild environment add REDUCE_WORKER_DEBUG true",
		Run:     runWrapper(runAddEnv),
	}
	cmdRmEnv = &cobra.Command{
		Use:     "remove NAME VALUE",
		Aliases: []string{"rm"},
		Short:   "Remove an environment variable",
		Long:    "Updates the environment in the current ACI's manifest to not contain the given variable",
		Example: "acbuild environment remove REDUCE_WORKER_DEBUG",
		Run:     runWrapper(runRemoveEnv),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdEnv)
	cmdEnv.AddCommand(cmdAddEnv)
	cmdEnv.AddCommand(cmdRmEnv)
}

func runAddEnv(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 2 {
		stderr("environment add: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Adding environment variable %q=%q", args[0], args[1])
	}

	err := lib.AddEnv(tmpacipath(), args[0], args[1])

	if err != nil {
		stderr("environment add: %v", err)
		return 1
	}

	return 0
}

func runRemoveEnv(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) > 1 {
		stderr("environment remove: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Removing environment variable %q", args[0])
	}

	err := lib.RemoveEnv(tmpacipath(), args[0])

	if err != nil {
		stderr("environment remove: %v", err)
		return 1
	}

	return 0
}
