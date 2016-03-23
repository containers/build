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
	cmdSetExec = &cobra.Command{
		Use:     "set-exec CMD [ARGS]",
		Short:   "Set the exec command",
		Long:    "Sets the exec command in the ACI's manifest",
		Example: "acbuild set-exec /usr/sbin/nginx -g \"daemon off;\"",
		Run:     runWrapper(runSetExec),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdSetExec)
}

func runSetExec(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Setting exec command %v", args)
	}

	err := newACBuild().SetExec(args)

	if err != nil {
		stderr("set-exec: %v", err)
		return getErrorCode(err)
	}

	return 0
}
