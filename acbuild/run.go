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

	"github.com/appc/acbuild/lib"
)

var (
	insecure = false
	cmdRun   = &cobra.Command{
		Use:     "run CMD [ARGS]",
		Short:   "Run a command in an ACI",
		Long:    "Run a given command in an ACI, and save the resulting container as a new ACI",
		Example: "acbuild run yum install nginx",
		Run:     runWrapper(runRun),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdRun)

	cmdRun.Flags().BoolVar(&insecure, "insecure", false, "Allows fetching dependencies over http")
}

func runRun(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}

	lockfile, err := getLock()
	if err != nil {
		stderr("run: %v", err)
		return 1
	}
	defer func() {
		if err := releaseLock(lockfile); err != nil {
			stderr("run: %v", err)
			exit = 1
		}
	}()

	if debug {
		stderr("Running: %v", args)
	}

	err = lib.Run(tmpacipath(), depstorepath(), targetpath(), scratchpath(), workpath(), args, insecure)

	if err != nil {
		stderr("run: %v", err)
		return 1
	}

	return 0
}
