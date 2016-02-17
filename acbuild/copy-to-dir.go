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

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	toDir        bool
	cmdCopyToDir = &cobra.Command{
		Use:     "copy-to-dir PATH1_ON_HOST PATH2_ON_HOST ... PATH_IN_ACI",
		Short:   "Copy a file or directory into a directory in an ACI",
		Example: "acbuild copy-to-dir build/bin/* /usr/bin",
		Run:     runWrapper(runCopyToDir),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdCopyToDir)
}

func runCopyToDir(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) < 2 {
		stderr("copy-to-dir: incorrect number of arguments")
		return 1
	}

	if debug {
		logMsg := "Copying "
		for i := 0; i < len(args)-1; i++ {
			logMsg += fmt.Sprintf("%s ", args[i])
		}
		logMsg += "to "
		logMsg += fmt.Sprintf("%s", args[len(args)-1])
		stderr(logMsg)
	}

	err := newACBuild().CopyToDir(args[:len(args)-1], args[len(args)-1])

	if err != nil {
		stderr("copy-to-dir: %v", err)
		return getErrorCode(err)
	}

	return 0
}
