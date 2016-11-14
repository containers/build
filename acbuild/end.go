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
	cmdEnd = &cobra.Command{
		Use:     "end",
		Short:   "end a current build",
		Long:    "End the current build, deleting the current context",
		Example: "acbuild end",
		Run:     runWrapper(runEnd),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdEnd)
}

func runEnd(cmd *cobra.Command, args []string) (exit int) {
	if len(args) != 0 {
		stderr("end: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Ending the build")
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.End()

	if err != nil {
		stderr("end: %v", err)
		return getErrorCode(err)
	}

	return 0
}
