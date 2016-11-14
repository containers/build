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
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	cmdIso = &cobra.Command{
		Use:   "isolator [command]",
		Short: "Manage isolators (appc only)",
	}
	cmdAddIso = &cobra.Command{
		Use:     "add NAME JSON_FILE",
		Short:   "Add an isolator (appc only)",
		Long:    "Updates the ACI to contain an isolator with the given name and value. If the isolator exists, its value will be changed.",
		Example: "acbuild isolator add resource/cpu ./value.json",
		Run:     runWrapper(runAddIso),
	}
	cmdRmIso = &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm"},
		Short:   "Remove an isolator (appc only)",
		Long:    "Updates the current ACI's manifest to not contain the given isolator",
		Example: "acbuild isolator remove resource/memory",
		Run:     runWrapper(runRemoveIso),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdIso)
	cmdIso.AddCommand(cmdAddIso)
	cmdIso.AddCommand(cmdRmIso)
}

func runAddIso(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 2 {
		stderr("isolator add: incorrect number of arguments")
		return 1
	}

	var r io.Reader
	var err error
	if args[1] == "-" {
		r = os.Stdin
	} else {
		file, err := os.Open(args[1])
		if err != nil {
			stderr("isolator add: %v", err)
			return 1
		}
		defer file.Close()
		r = file
	}

	val, err := ioutil.ReadAll(r)
	if err != nil {
		stderr("isolator add: %v", err)
		return 1
	}

	if debug {
		stderr("Adding isolator %q=%q", args[0], string(val))
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.AddIsolator(args[0], val)

	if err != nil {
		stderr("isolator add: %v", err)
		return 1
	}

	return 0
}

func runRemoveIso(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) > 1 {
		stderr("isolator remove: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Removing isolator %q", args[0])
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.RemoveIsolator(args[0])

	if err != nil {
		stderr("isolator remove: %v", err)
		return 1
	}

	return 0
}
