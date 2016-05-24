// Copyright 2016 The rkt Authors
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

// +build manpages

package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	cmdManPages = &cobra.Command{
		Use:     "gen-man-pages",
		Short:   "Generate man pages for acbuild",
		Example: "acbuild gen-man-pages",
		Run:     runWrapper(runManPages),
		Hidden:  true,
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdManPages)
}

func runManPages(cmd *cobra.Command, args []string) (exit int) {
	outputDir := "./man/"

	os.MkdirAll(outputDir, 0755)

	header := &doc.GenManHeader{
		Title: "MINE",
		//Section string
		//Date    *time.Time
		//date    string
		//Source  string
		//Manual  string
	}
	err := doc.GenManTree(cmdAcbuild, header, outputDir)
	if err != nil {
		stderr("gen-man-pages: %v", err)
		return 1
	}

	stderr("man pages have been output to %s", outputDir)
	return 0
}
