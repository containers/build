// Copyright 2016 The appc Authors
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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	nestedScriptEnvVar = "ACBUILD_NESTED_SCRIPT"
)

var (
	errSingleQuote = fmt.Errorf("unterminated single quote block")
	errDoubleQuote = fmt.Errorf("unterminated double quote block")
	errEscape      = fmt.Errorf("ended with an escape")
	cmdScript      = &cobra.Command{
		Use:     "script SCRIPT_FILE",
		Short:   "Runs an acbuild script",
		Example: "acbuild script build-myapp.acb",
		Run:     runWrapper(runScript),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdScript)
}

func runScript(cmd *cobra.Command, args []string) (exit int) {
	if len(args) != 1 {
		cmd.Usage()
		return 1
	}

	scriptName := args[0]
	rawScript, err := ioutil.ReadFile(scriptName)
	if err != nil {
		stderr("script: %v", err)
		return getErrorCode(err)
	}

	if debug {
		stderr("Running script from %s", scriptName)
	}

	err = execScript(rawScript)
	if err != nil {
		stderr("script: %v", err)
		return getErrorCode(err)
	}
	return 0
}

func execScript(rawScript []byte) error {
	script := strings.Split(string(rawScript), "\n")
	for i, s := range script {
		s = strings.TrimSpace(s)
		script[i] = s

		if strings.HasPrefix(strings.ToLower(s), "run") && os.Geteuid() != 0 {
			return fmt.Errorf("scripts using the run subcommand must be run as root")
		}
	}
	script = joinLines(script)

	var tmpDir string
	nestedScript := false
	if inheritedPath := os.Getenv(nestedScriptEnvVar); inheritedPath != "" {
		tmpDir = inheritedPath
		nestedScript = true
	} else {
		var err error
		tmpDir, err = ioutil.TempDir("", "acbuild")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)
		contextpath = tmpDir
	}

	for _, line := range script {
		if line == "" {
			continue
		}
		err := execACBuild(tmpDir, line)
		if err != nil {
			if !strings.HasPrefix(line, "begin") && !nestedScript {
				err1 := newACBuild().End()
				if err1 != nil {
					stderr("script: %v", err1)
				}
			}
			return err
		}
	}
	if !nestedScript {
		err := newACBuild().End()
		if err != nil {
			return err
		}
	}
	if debug {
		stderr("Script has been completed")
	}

	return nil
}

func execACBuild(workPath, line string) error {
	suppliedArgs, err := tokenizeLine(line)
	if err != nil {
		return err
	}
	if len(suppliedArgs) == 0 {
		return nil
	}
	suppliedArgs[0] = strings.ToLower(suppliedArgs[0])
	if suppliedArgs[0] == "run" || suppliedArgs[0] == "set-exec" {
		suppliedArgs = insertRunTacks(suppliedArgs)
	}
	args := []string{"--debug", "--work-path=" + workPath}
	args = append(args, suppliedArgs...)
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()
	if suppliedArgs[0] == "script" {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", nestedScriptEnvVar, workPath))
	}
	return cmd.Run()
}

func joinLines(script []string) []string {
	for i, line := range script {
		if strings.HasSuffix(line, `\`) && i != len(script)-1 {
			line = strings.TrimSuffix(line, `\`)
			script[i] = ""
			script[i+1] = line + " " + script[i+1]
		}
	}
	return script
}

func tokenizeLine(line string) ([]string, error) {
	var tokens []string
	buf := &bytes.Buffer{}
	inSingleQuoteBlock := false
	inDoubleQuoteBlock := false
	isEscaped := false
lineLoop:
	for _, char := range line {
		if isEscaped {
			buf.WriteRune(char)
			isEscaped = false
			continue
		}
		switch {
		case char == '\\':
			isEscaped = true
		case (char == ' ' || char == '	') && !inSingleQuoteBlock && !inDoubleQuoteBlock:
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
		case char == '\'' && !inDoubleQuoteBlock:
			inSingleQuoteBlock = !inSingleQuoteBlock
		case char == '"' && !inSingleQuoteBlock:
			inDoubleQuoteBlock = !inDoubleQuoteBlock
		case char == '#' && !inSingleQuoteBlock && !inDoubleQuoteBlock:
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			break lineLoop
		default:
			buf.WriteRune(char)
		}
	}
	if inSingleQuoteBlock {
		return nil, errSingleQuote
	}
	if inDoubleQuoteBlock {
		return nil, errDoubleQuote
	}
	if isEscaped {
		return nil, errEscape
	}
	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}
	return tokens, nil
}

func insertRunTacks(tokens []string) []string {
	insertLocation := -1
	for i, tok := range tokens {
		if i == 0 {
			continue
		}
		if tok == "--" {
			return tokens
		}
		if tok[0] != '-' {
			insertLocation = i
			break
		}
	}
	if insertLocation == -1 {
		return tokens
	}
	newTokens := make([]string, insertLocation)
	copy(newTokens, tokens[:insertLocation])
	newTokens = append(newTokens, "--")
	newTokens = append(newTokens, tokens[insertLocation:]...)
	return newTokens
}
