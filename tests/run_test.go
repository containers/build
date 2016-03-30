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
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"
)

const goprogram = `
package main

import (
	"fmt"
)

func main() {
	fmt.Printf("success")
}
`

func TestRun(t *testing.T) {
	if os.Getenv("ENABLE_SYSTEMD_TESTS") == "" {
		t.Skip("skipping test; $ENABLE_SYSTEMD_TESTS not set")
	}

	// Build a statically linked test program
	tmpsourcedir := mustTempDir()
	defer os.RemoveAll(tmpsourcedir)
	tmpsource := path.Join(tmpsourcedir, "thing.go")
	defer os.RemoveAll(tmpsourcedir)
	err := ioutil.WriteFile(tmpsource, []byte(goprogram), 0644)
	if err != nil {
		panic(err)
	}

	tmprootfs := mustTempDir()
	defer os.RemoveAll(tmprootfs)

	cmd := exec.Command("go", "build", "-o", path.Join(tmprootfs, "worker"), "-tags", "netgo", "-ldflags", "-w", tmpsource)
	cmd.Env = []string{"CGO_ENABLED=0", "GOOS=linux", "GOROOT=" + os.Getenv("GOROOT"), "GOPATH=" + os.Getenv("GOPATH")}
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		panic(err)
	}

	// Call begin on it
	tmpdir := mustTempDir()
	_, _, _, err = runACBuild(tmpdir, "begin", tmprootfs)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer os.RemoveAll(tmpdir)

	// acbuild run the binary
	_, stdout, stderr, err := runACBuild(tmpdir, "--no-history", "run", "/worker")
	if err != nil {
		panic(err)
	}
	if stderr != "" {
		t.Errorf("stderr wasn't empty: %s", stderr)
	}
	if stdout != "success" {
		t.Errorf("unexpected stdout: %s", stdout)
	}
}

func TestRunBadEngine(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	_, stdout, stderr, err := runACBuild(workingDir, "run", "--engine=invalid-engine", "command")
	if err == nil {
		t.Errorf("was not expecting err to be nil when run with invalid engine")
	}

	if stdout != "" {
		t.Errorf("printed to stdout when should not have: %s", stdout)
	}

	if stderr != "run: no such engine \"invalid-engine\"\n" {
		t.Errorf("unexpected message on stderr: %s", stderr)
	}
}
