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
	overwrite = false
	sign      = false
	cmdWrite  = &cobra.Command{
		Use:     "write ACI_PATH",
		Short:   "Write the image from the current build to a file",
		Example: "acbuild write --sign mynewapp.aci -- --no-default-keyring --keyring ./rkt.gpg",
		Run:     runWrapper(runWrite),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdWrite)

	cmdWrite.Flags().BoolVar(&overwrite, "overwrite", false, "overwrite the resulting ACI")
	cmdWrite.Flags().BoolVar(&sign, "sign", false, "(removed) sign the resulting ACI")
}

func runWrite(cmd *cobra.Command, args []string) (exit int) {
	if sign {
		stderr("write: the sign flag has been removed, please invoke gpg directly")
		return 1
	}

	if len(args) != 1 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Writing ACI to %s", args[0])
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.Write(args[0], overwrite)

	if err != nil {
		stderr("write: %v", err)
		return getErrorCode(err)
	}

	return 0
}
