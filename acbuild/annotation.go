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
	cmdAnno = &cobra.Command{
		Use:   "annotation [command]",
		Short: "Manage annotations",
	}
	cmdAddAnno = &cobra.Command{
		Use:     "add NAME VALUE",
		Short:   "Add an annotation",
		Long:    "Updates the ACI to contain an annotation with the given name and value. If the annotation already exists, its value will be changed.",
		Example: "acbuild annotation add documentation https://example.com/docs",
		Run:     runWrapper(runAddAnno),
	}
	cmdRmAnno = &cobra.Command{
		Use:     "remove NAME",
		Aliases: []string{"rm"},
		Short:   "Remove an annotation",
		Long:    "Removes the annotation with the given name from the ACI's manifest",
		Example: "acbuild annotation remove documentation",
		Run:     runWrapper(runRmAnno),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdAnno)
	cmdAnno.AddCommand(cmdAddAnno)
	cmdAnno.AddCommand(cmdRmAnno)
}

func runAddAnno(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 2 {
		stderr("annotation add: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Adding annotation %q=%q", args[0], args[1])
	}

	err := newACBuild().AddAnnotation(args[0], args[1])

	if err != nil {
		stderr("annotation add: %v", err)
		return getErrorCode(err)
	}

	return 0
}

func runRmAnno(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("annotation remove: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Removing annotation %q", args[0])
	}

	err := newACBuild().RemoveAnnotation(args[0])

	if err != nil {
		stderr("annotation remove: %v", err)
		return getErrorCode(err)
	}

	return 0
}
