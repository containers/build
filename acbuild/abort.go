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
	"path"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"

	"github.com/appc/acbuild/lib"
)

var (
	cmdAbort = &cobra.Command{
		Use:     "abort",
		Short:   "Abort an existing build",
		Long:    "Abort the current build, throwing away any changes since init was called",
		Example: "acbuild abort",
		Run:     runWrapper(runAbort),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdAbort)
}

func runAbort(cmd *cobra.Command, args []string) (exit int) {
	if len(args) != 0 {
		stderr("abort: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Aborting the build")
	}

	err := lib.Abort(path.Join(contextpath, workprefix))

	if err != nil {
		stderr("abort: %v", err)
		return 1
	}

	return 0
}
