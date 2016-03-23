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
	"github.com/spf13/cobra"
)

var (
	cmdSetUser = &cobra.Command{
		Use:     "set-user USER",
		Short:   "Set the user",
		Long:    "Set the user the app will run as inside the container",
		Example: "acbuild set-user www-data",
		Run:     runWrapper(runSetUser),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdSetUser)
}

func runSetUser(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("set-user: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Setting user to %s", args[0])
	}

	err := newACBuild().SetUser(args[0])

	if err != nil {
		stderr("set-user: %v", err)
		return getErrorCode(err)
	}

	return 0
}
