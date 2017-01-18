// Copyright 2015 The rkt Authors
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
	"crypto/sha512"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/tabwriter"
	"text/template"

	"github.com/coreos/rkt/pkg/multicall"
	"github.com/spf13/cobra"

	"github.com/containers/build/lib"
	"github.com/containers/build/lib/appc"
)

const (
	cliName = "acbuild"

	commandUsage = `\
NAME:
{{printf "\t%s - %s" .Name .Short}}

USAGE:
{{printf "\t%s" .UseLine}}

{{if (ne .Example "")}}\
EXAMPLE:
{{printf "\t%s" .Example}}

{{end}}\
\
{{if (ne (len .Commands) 0)}}\
COMMANDS:
{{range .Commands}}\
{{if (ne (len .Commands) 0)}}\
{{printf "\t%s%s\t%s" .Name (subcmdList .Commands) .Short}}
{{else}}\
{{printf "\t%s\t%s" .Name .Short}}
{{end}}\
{{end}}\

{{end}}\
\
OPTIONS:
{{.LocalFlags.FlagUsages}}\
{{if (eq .Name "acbuild")}}\

DOCUMENTATION:
	Additional documentation is available at https://github.com/containers/build\
{{end}}
`
)

var (
	debug          bool
	contextpath    string
	aciToModify    string
	ociToModify    string
	disableHistory bool

	cmdExitCode int

	errCobra = fmt.Errorf("cobra error")

	templFuncs = template.FuncMap{
		"subcmdList": func(cmds []*cobra.Command) string {
			var subcmds []string
			for _, subcmd := range cmds {
				subcmds = append(subcmds, subcmd.Name())
			}
			return " [" + strings.Join(subcmds, "|") + "]"
		},
	}

	commandUsageTemplate = template.Must(template.New("command_usage").Funcs(templFuncs).Parse(strings.Replace(commandUsage, "\\\n", "", -1)))
)

var cmdAcbuild = &cobra.Command{
	Use:   "acbuild [command]",
	Short: "the application container build system",
}

func init() {
	cmdAcbuild.PersistentFlags().BoolVar(&debug, "debug", false, "Print out debug information to stderr")
	cmdAcbuild.PersistentFlags().StringVar(&contextpath, "work-path", ".", "Path to place working files in")
	cmdAcbuild.PersistentFlags().StringVar(&aciToModify, "modify-appc", "", "Path to an ACI to modify (ignores build context)")
	cmdAcbuild.PersistentFlags().StringVar(&ociToModify, "modify-oci", "", "Path to an OCI image to modify (ignores build context)")
	cmdAcbuild.PersistentFlags().BoolVar(&disableHistory, "no-history", false, "Don't add annotations with the command that was run")

	cobra.EnablePrefixMatching = true
}

func newACBuild() (*lib.ACBuild, error) {
	bmode, err := lib.GetBuildMode(contextpath)
	if err != nil {
		return nil, err
	}
	return lib.NewACBuild(contextpath, debug, bmode)
}

func newACBuildWithBuildMode(bmode lib.BuildMode) (*lib.ACBuild, error) {
	return lib.NewACBuild(contextpath, debug, bmode)
}

func getErrorCode(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.Sys().(syscall.WaitStatus).ExitStatus()
	}
	switch err {
	case appc.ErrNotFound:
		return 2
	case errCobra:
		return 3
	case nil:
		return 0
	default:
		return 1
	}
}

// runWrapper return a func(cmd *cobra.Command, args []string) that internally
// will add command function return code and the reinsertion of the "--" flag
// terminator.
func runWrapper(cf func(cmd *cobra.Command, args []string) (exit int)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if aciToModify == "" && ociToModify == "" {
			cmdExitCode = cf(cmd, args)
			switch cmd.Name() {
			case "cat-manifest", "begin", "write", "end", "version", "gen-man-pages", "script":
				return
			}
			if cmdExitCode == 0 && !disableHistory {
				err := addACBuildAnnotation(cmd, args)
				if err != nil {
					stderr("%v", err)
					cmdExitCode = 1
					return
				}
			}
			return
		}

		if aciToModify != "" && ociToModify != "" {
			stderr("can't modify an appc image and an oci image at the same time")
			cmdExitCode = 1
			return
		}

		switch cmd.Name() {
		case "begin", "write", "end", "version", "gen-man-pages", "script":
			stderr("Can't use --modify flags with %s.", cmd.Name())
			cmdExitCode = 1
			return
		}

		toModify := aciToModify
		if ociToModify != "" {
			toModify = ociToModify
		}

		finfo, err := os.Stat(toModify)
		switch {
		case os.IsNotExist(err):
			stderr("image doesn't appear to exist: %s.", toModify)
			cmdExitCode = 1
			return
		case err != nil:
			stderr("error accessing image to modify: %v.", err)
			cmdExitCode = 1
			return
		case finfo.IsDir():
			stderr("image to modify is a directory: %s.", toModify)
			cmdExitCode = 1
			return
		}

		absoluteToModify, err := filepath.Abs(toModify)
		if err != nil {
			stderr("%v", err)
			cmdExitCode = 1
			return
		}

		hash := sha512.New().Sum([]byte(absoluteToModify))
		contextpath := path.Join(os.TempDir(), fmt.Sprintf("acbuild-%x", hash))

		if len(contextpath) > 16 {
			contextpath = contextpath[:16]
		}

		err = os.MkdirAll(contextpath, 0755)
		if err != nil {
			stderr("%v", err)
			cmdExitCode = 1
			return
		}
		defer os.RemoveAll(contextpath)

		modifyMode := lib.BuildModeAppC
		if ociToModify != "" {
			modifyMode = lib.BuildModeOCI
		}

		a, err := newACBuildWithBuildMode(modifyMode)
		if err != nil {
			stderr("%v", err)
			cmdExitCode = 1
			return
		}

		err = a.Begin(absoluteToModify, false, modifyMode)
		if err != nil {
			stderr("%v", err)
			cmdExitCode = getErrorCode(err)
			return
		}

		defer func() {
			err = a.End()
			if err != nil {
				stderr("%v", err)
				if cmdExitCode == 0 {
					cmdExitCode = getErrorCode(err)
				}
			}
		}()

		cmdExitCode = cf(cmd, args)

		if cmdExitCode == 0 && !disableHistory {
			err := addACBuildAnnotation(cmd, args)
			if err != nil {
				stderr("%v", err)
				cmdExitCode = 1
				return
			}
		}

		dir, file := path.Split(aciToModify)
		tmpFile := path.Join(dir, "."+file+".tmp")

		err = a.Write(tmpFile, true)
		if err != nil {
			stderr("%v", err)
			cmdExitCode = getErrorCode(err)
			return
		}

		err = os.Rename(tmpFile, toModify)
		if err != nil {
			os.Remove(tmpFile)
			stderr("%v", err)
			cmdExitCode = getErrorCode(err)
			return
		}

	}
}

func main() {
	multicall.Add("acbuild-script", func() error {
		cmd := exec.Command("acbuild", append([]string{"script"}, os.Args[1:]...)...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	})
	// check if acbuild is executed with a multicall command
	multicall.MaybeExec()

	cmdAcbuild.SetUsageFunc(func(cmd *cobra.Command) error {
		tabOut := new(tabwriter.Writer)
		tabOut.Init(os.Stdout, 0, 8, 1, '\t', 0)
		commandUsageTemplate.Execute(tabOut, cmd)
		tabOut.Flush()
		return nil
	})

	// Make help just show the usage
	cmdAcbuild.SetHelpTemplate(`{{.UsageString}}`)

	err := cmdAcbuild.Execute()
	if cmdExitCode == 0 && err != nil {
		cmdExitCode = getErrorCode(errCobra)
	}
	os.Exit(cmdExitCode)
}

func stderr(format string, a ...interface{}) {
	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stderr, strings.TrimSuffix(out, "\n"))
}

func stdout(format string, a ...interface{}) {
	out := fmt.Sprintf(format, a...)
	fmt.Fprintln(os.Stdout, strings.TrimSuffix(out, "\n"))
}

func addACBuildAnnotation(cmd *cobra.Command, args []string) error {
	const annoNamePattern = "coreos.com/acbuild/command-%d"

	acb, err := newACBuild()
	if err != nil {
		return err
	}

	annotations, err := acb.GetAnnotations()
	if err != nil {
		return err
	}

	var acbuildCount int
	for name, _ := range annotations {
		var tmpCount int
		n, _ := fmt.Sscanf(string(name), annoNamePattern, &tmpCount)
		if n == 1 && tmpCount > acbuildCount {
			acbuildCount = tmpCount
		}
	}

	command := cmd.Name()
	tmpcmd := cmd.Parent()
	for {
		command = tmpcmd.Name() + " " + command
		if tmpcmd == cmdAcbuild {
			break
		}
		tmpcmd = tmpcmd.Parent()
	}

	for _, a := range args {
		command += fmt.Sprintf(" %q", a)
	}

	return acb.AddAnnotation(fmt.Sprintf(annoNamePattern, acbuildCount+1), command)
}
