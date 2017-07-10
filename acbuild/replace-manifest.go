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
	cmdReplaceMan = &cobra.Command{
		Use:     "replace-manifest PATH_TO_MANIFEST",
		Short:   "Replace the manifest in the current build",
		Example: "acbuild replace-manifest path/to/manifest",
		Run:     runWrapper(runReplaceMan),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdReplaceMan)
}

func runReplaceMan(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("replace-manifest: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Replacing manifest in ACI with the manifest at ", args[0])
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.Replace(args[0])

	if err != nil {
		stderr("replace-manifest: %v", err)
		return 1
	}

	return 0
}
