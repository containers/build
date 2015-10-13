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
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
	"text/tabwriter"
	"text/template"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/multicall"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/pflag"

	"github.com/appc/acbuild/util"
)

const (
	cliName        = "acbuild"
	cliDescription = "acbuild, the application container build system"
)

var (
	workprefix = ".acbuild"

	commandUsageTemplate *template.Template

	templFuncs = template.FuncMap{
		"cmdName": func(cmd *cobra.Command, startCmd *cobra.Command) string {
			parts := []string{cmd.Name()}
			for cmd.HasParent() && cmd.Parent().Name() != startCmd.Name() {
				cmd = cmd.Parent()
				parts = append([]string{cmd.Name()}, parts...)
			}
			return strings.Join(parts, " ")
		},
	}

	debug       bool
	contextpath string

	tabOut      *tabwriter.Writer
	cmdExitCode int
)

func tmpacipath() string {
	return path.Join(contextpath, workprefix, "currentaci")
}

func targetpath() string {
	return path.Join(contextpath, workprefix, "target")
}

func scratchpath() string {
	return path.Join(contextpath, workprefix, "scratch")
}

func depstorepath() string {
	return path.Join(contextpath, workprefix, "depstore")
}

func workpath() string {
	return path.Join(contextpath, workprefix, "work")
}

func lockpath() string {
	return path.Join(contextpath, workprefix, "lock")
}

var cmdAcbuild = &cobra.Command{
	Use:   "acbuild [command]",
	Short: cliDescription,
}

func init() {
	cmdAcbuild.PersistentFlags().BoolVar(&debug, "debug", false, "Print out debug information to stderr")
	cmdAcbuild.PersistentFlags().StringVar(&contextpath, "work-path", ".", "Path to place working files in")
}

func init() {
	tabOut = new(tabwriter.Writer)
	tabOut.Init(os.Stdout, 0, 8, 1, '\t', 0)

	cobra.EnablePrefixMatching = true

	commandUsage := `
{{ $cmd := .Cmd }}\
{{ $cmdname := cmdName .Cmd .Cmd.Root }}\
NAME:
{{ if not .Cmd.HasParent }}\
{{printf "\t%s - %s" .Cmd.Name .Cmd.Short}}
{{else}}\
{{printf "\t%s - %s" $cmdname .Cmd.Short}}
{{end}}\

USAGE:
{{printf "\t%s" .Cmd.UseLine}}
{{if .Cmd.HasSubCommands}}\

COMMANDS:
{{range .SubCommands}}\
{{ $cmdname := cmdName . $cmd }}\
{{ if .Runnable }}\
{{printf "\t%s\t%s" $cmdname .Short}}
{{end}}\
{{end}}\
{{end}}\
{{if .Cmd.HasLocalFlags}}\

OPTIONS:
{{.Cmd.LocalFlags.FlagUsages}}\
{{end}}\
{{if .Cmd.HasInheritedFlags}}\

GLOBAL OPTIONS:
{{.Cmd.InheritedFlags.FlagUsages}}\
{{end}}
`[1:]

	commandUsageTemplate = template.Must(template.New("command_usage").Funcs(templFuncs).Parse(strings.Replace(commandUsage, "\\\n", "", -1)))
}

// runWrapper return a func(cmd *cobra.Command, args []string) that internally
// will add command function return code and the reinsertion of the "--" flag
// terminator.
func runWrapper(cf func(cmd *cobra.Command, args []string) (exit int)) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		cmdExitCode = cf(cmd, args)
	}
}

func main() {
	// check if acbuild is executed with a multicall command
	multicall.MaybeExec()

	cmdAcbuild.SetUsageFunc(usageFunc)

	// Make help just show the usage
	cmdAcbuild.SetHelpTemplate(`{{.UsageString}}`)

	cmdAcbuild.Execute()
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

func getLock() (*os.File, error) {
	ex, err := util.Exists(path.Join(contextpath, workprefix))
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, fmt.Errorf("build not in progress in this working dir - try \"acbuild begin\"")
	}

	lockfile, err := os.OpenFile(lockpath(), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(lockfile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		if err == syscall.EWOULDBLOCK {
			return nil, fmt.Errorf("lock already held - is another acbuild running in this working dir?")
		}
		return nil, err
	}

	return lockfile, nil
}

func releaseLock(lockfile *os.File) error {
	err := syscall.Flock(int(lockfile.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}

	err = lockfile.Close()
	if err != nil {
		return err
	}
	lockfile = nil

	err = os.Remove(lockpath())
	if err != nil {
		return err
	}

	return nil
}

func getSubCommands(cmd *cobra.Command) []*cobra.Command {
	subCommands := []*cobra.Command{}
	for _, subCmd := range cmd.Commands() {
		subCommands = append(subCommands, subCmd)
		subCommands = append(subCommands, getSubCommands(subCmd)...)
	}
	return subCommands
}

func usageFunc(cmd *cobra.Command) error {
	subCommands := getSubCommands(cmd)
	commandUsageTemplate.Execute(tabOut, struct {
		Executable  string
		Cmd         *cobra.Command
		CmdFlags    *pflag.FlagSet
		SubCommands []*cobra.Command
	}{
		cliName,
		cmd,
		cmd.Flags(),
		subCommands,
	})
	tabOut.Flush()
	return nil
}
