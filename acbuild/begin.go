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
	"os"
	"path"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"

	"github.com/appc/acbuild/lib"
	"github.com/appc/acbuild/util"
)

var (
	cmdBegin = &cobra.Command{
		Use:     "begin [START_ACI_PATH]",
		Short:   "Start a new build",
		Long:    "Begins a new build. By default operations will be performed on top of an empty image, but a start image can be provided",
		Example: "acbuild begin",
		Run:     runWrapper(runBegin),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdBegin)
}

func runBegin(cmd *cobra.Command, args []string) (exit int) {
	if len(args) > 1 {
		stderr("begin: incorrect number of arguments")
		return 1
	}

	ex, err := util.Exists(path.Join(contextpath, workprefix))
	if err != nil {
		stderr("begin: %v", err)
		return 1
	}
	if ex {
		stderr("begin: build already in progress in this working dir")
		return 1
	}

	err = os.MkdirAll(path.Join(contextpath, workprefix), 0755)
	if err != nil {
		stderr("begin: %v", err)
		return 1
	}

	lockfile, err := getLock()
	if err != nil {
		stderr("begin: %v", err)
		return 1
	}
	defer func() {
		if err := releaseLock(lockfile); err != nil {
			stderr("begin: %v", err)
			exit = 1
		}
	}()

	if debug {
		if len(args) == 0 {
			stderr("Beginning build with an empty ACI")
		} else {
			stderr("Beginning build with %s", args[0])
		}
	}

	if len(args) == 0 {
		err = lib.Begin(tmpacipath(), "")
	} else {
		err = lib.Begin(tmpacipath(), args[0])
	}

	if err != nil {
		stderr("begin: %v", err)
		return 1
	}

	return 0
}
