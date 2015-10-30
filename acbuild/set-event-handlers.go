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
	cmdSetEventHandler = &cobra.Command{
		Use:     "set-event-handler [command]",
		Aliases: []string{"set-eh"},
		Short:   "Manage event handlers",
	}
	cmdSetPreStart = &cobra.Command{
		Use:     "pre-start CMD [ARGS]",
		Short:   "Set the pre-start event handler",
		Example: "acbuild set-event-handler pre-start /root/setup-stuff.sh",
		Run:     runWrapper(runSetPreStart),
	}
	cmdSetPostStop = &cobra.Command{
		Use:     "post-stop CMD [ARGS]",
		Short:   "Set the post-stop event handler",
		Example: "acbuild set-event-handler post-stop /bin/report-results.sh",
		Run:     runWrapper(runSetPostStop),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdSetEventHandler)
	cmdSetEventHandler.AddCommand(cmdSetPreStart)
	cmdSetEventHandler.AddCommand(cmdSetPostStop)
}

func runSetPreStart(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Setting pre-start event handler to %v", args)
	}

	err := newACBuild().SetPreStart(args)

	if err != nil {
		stderr("pre-start: %v", err)
		return 1
	}

	return 0
}

func runSetPostStop(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}

	if debug {
		stderr("Setting post-stop event handler to %v", args)
	}

	err := newACBuild().SetPostStop(args)

	if err != nil {
		stderr("post-stop: %v", err)
		return 1
	}

	return 0
}
