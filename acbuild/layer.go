// Copyright 2017 The acbuild Authors
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
	cmdLayer = &cobra.Command{
		Use:     "layer",
		Short:   "Creates a new layer in the image (OCI only)",
		Example: "acbuild layer",
		Run:     runWrapper(runLayer),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdLayer)
}

func runLayer(cmd *cobra.Command, args []string) (exit int) {
	if len(args) != 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Adding new layer")
	}

	a, err := newACBuild()
	if err != nil {
		stderr("%v", err)
		return 1
	}
	err = a.NewLayer()

	if err != nil {
		stderr("layer: %v", err)
		return getErrorCode(err)
	}

	return 0
}
