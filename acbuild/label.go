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
	cmdLabel = &cobra.Command{
		Use:   "label [command]",
		Short: "Manage labels",
	}
	cmdAddLabel = &cobra.Command{
		Use:     "add NAME VALUE",
		Short:   "Add a label, or update an existing one",
		Example: "acbuild label add arch amd64",
		Run:     runWrapper(runAddLabel),
	}
	cmdRmLabel = &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm"},
		Short:   "Remove a label",
		Example: "acbuild label remove arch",
		Run:     runWrapper(runRemoveLabel),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdLabel)
	cmdLabel.AddCommand(cmdAddLabel)
	cmdLabel.AddCommand(cmdRmLabel)
}

func runAddLabel(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 2 {
		stderr("label add: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Adding label %q=%q", args[0], args[1])
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.AddLabel(args[0], args[1])

	if err != nil {
		stderr("label add: %v", err)
		return getErrorCode(err)
	}

	return 0
}

func runRemoveLabel(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("label remove: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Removing label %q", args[0])
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.RemoveLabel(args[0])

	if err != nil {
		stderr("label remove: %v", err)
		return getErrorCode(err)
	}

	return 0
}
