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
	"fmt"

	"github.com/spf13/cobra"
)

var (
	prettyPrint bool
	printFile   string

	cmdCat = &cobra.Command{
		Use:     "cat-manifest",
		Short:   "Print the manifest from the current build",
		Example: "acbuild cat-manifest",
		Run:     runWrapper(runCat),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdCat)
	cmdCat.Flags().BoolVar(&prettyPrint, "pretty-print", false, "Print the manifest with whitespace")
	cmdCat.Flags().StringVar(&printFile, "file", "manifest", "(oci only) Which file to print, accepts \"manifest\" (the default), and \"config\"")
}

func runCat(cmd *cobra.Command, args []string) (exit int) {
	if len(args) != 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Printing manifest from current build")
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	switch printFile {
	case "manifest":
		err = a.Print(prettyPrint, false)
	case "config":
		err = a.Print(prettyPrint, true)
	default:
		err = fmt.Errorf("don't know how to print a %q", printFile)
	}
	if err != nil {
		stderr("cat-manifest: %v", err)
		return 1
	}

	return 0
}
