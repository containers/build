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
	prettyPrint bool
	printConfig bool

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
	cmdCat.Flags().BoolVar(&printConfig, "config", false, "(oci only) Prints the image config instead of the manifest")
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
	err = a.Print(prettyPrint, printConfig)
	if err != nil {
		stderr("cat-manifest: %v", err)
		return 1
	}

	return 0
}

func runCatOnACI(aciToModify string) int {
	// TODO
	//err := lib.Print(aciToModify, prettyPrint)
	//if err != nil {
	//	stderr("cat-manifest: %v", err)
	//	return 1
	//}
	return 0
}
