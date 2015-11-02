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

package tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"syscall"
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/kylelemons/godebug/pretty"
)

var (
	acbuildBinPath string

	systemLabels = types.Labels{
		types.Label{
			*types.MustACIdentifier("arch"),
			runtime.GOARCH,
		},
		types.Label{
			*types.MustACIdentifier("os"),
			runtime.GOOS,
		},
	}

	emptyManifest = schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:  nil,
			User:  "0",
			Group: "0",
		},
		Labels: systemLabels,
	}
)

func init() {
	acbuildBinPath = os.Getenv("ACBUILD_BIN")
	if acbuildBinPath == "" {
		fmt.Fprintf(os.Stderr, "ACBUILD_BIN environmment variable must be set\n")
		os.Exit(1)
	} else if _, err := os.Stat(acbuildBinPath); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type acbuildError struct {
	err      *exec.ExitError
	exitCode int
	stdout   []byte
	stderr   []byte
}

func (ae acbuildError) Error() string {
	return fmt.Sprintf("non-zero exit code of %d: %v\nstdout:\n%s\nstderr:\n%s", ae.exitCode, ae.err, string(ae.stdout), string(ae.stderr))
}

func runACBuild(workingDir string, args ...string) *acbuildError {
	cmd := exec.Command(acbuildBinPath, args...)
	cmd.Dir = workingDir
	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stderrpipe, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	readTillClosed := func(in io.ReadCloser) []byte {
		msg, err := ioutil.ReadAll(in)
		if err != nil {
			panic(err)
		}
		return msg
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	stdout := readTillClosed(stdoutpipe)
	stderr := readTillClosed(stderrpipe)
	err = cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		return &acbuildError{exitErr, code, stdout, stderr}
	}
	if err != nil {
		panic(err)
	}
	return nil
}

func setUpTest(t *testing.T) string {
	tmpdir := mustTempDir()

	err := runACBuild(tmpdir, "begin")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	return tmpdir
}

func cleanUpTest(tmpdir string) error {
	return os.RemoveAll(tmpdir)
}

func mustTempDir() string {
	dir, err := ioutil.TempDir("", "acbuild-test")
	if err != nil {
		panic(err)
	}
	return dir
}

func checkManifest(t *testing.T, workingDir string, wantedManifest schema.ImageManifest) {
	acipath := path.Join(workingDir, ".acbuild", "currentaci")

	manblob, err := ioutil.ReadFile(path.Join(acipath, aci.ManifestFile))
	if err != nil {
		panic(err)
	}

	var man schema.ImageManifest

	err = man.UnmarshalJSON(manblob)
	if err != nil {
		t.Errorf("invalid manifest schema: %v", err)
	}

	if str := pretty.Compare(man, wantedManifest); str != "" {
		t.Errorf("unexpected manifest:\n%s", str)
	}
}

func checkEmptyRootfs(t *testing.T, workingDir string) {
	files, err := ioutil.ReadDir(path.Join(workingDir, ".acbuild", "currentaci", aci.RootfsDir))
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(files) != 0 {
		t.Errorf("rootfs in aci contains files, should be empty")
	}
}
