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
	explicitTarget bool
	cmdCopy        = &cobra.Command{
		Use:     "copy PATH_ON_HOST... PATH_IN_ACI",
		Short:   "Copy a file or directory into an ACI",
		Example: "acbuild copy stuff/* nginx.conf /etc/nginx/",
		Run:     runWrapper(runCopy),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdCopy)
	cmdCopy.Flags().BoolVarP(&explicitTarget, "explicit-target", "T", false, "copy a single file/directory to the specified path")
}

func runCopy(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) < 2 || (explicitTarget && len(args) != 2) {
		stderr("copy: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Copying host:%s to aci:%s", args[0], args[1])
	}

	var err error
	if explicitTarget {
		err = newACBuild().CopyToTarget(args[0], args[1])
	} else {
		err = newACBuild().CopyToDir(args[:len(args)-1], args[len(args)-1])
	}

	if err != nil {
		stderr("copy: %v", err)
		return getErrorCode(err)
	}

	return 0
}
