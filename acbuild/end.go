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
	overwrite = false
	cmdEnd    = &cobra.Command{
		Use:     "end ACI_PATH",
		Short:   "End a build",
		Long:    "Ends a running build, placing the resulting ACI at the provided path",
		Example: "acbuild end mynewapp.aci",
		Run:     runWrapper(runEnd),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdEnd)

	cmdEnd.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite output ACI")
}

func runEnd(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) > 1 {
		stderr("end: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Ending build. Writing completed ACI to %s", args[0])
	}

	err := lib.End(tmpacipath(), args[0], path.Join(contextpath, workprefix), overwrite)

	if err != nil {
		stderr("end: %v", err)
		return 1
	}

	return 0
}
