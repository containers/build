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
	"bytes"
	"fmt"
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

func emptyManifest() schema.ImageManifest {
	return schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		Labels:    systemLabels,
	}
}

func emptyManifestWithApp() schema.ImageManifest {
	return schema.ImageManifest{
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
}

func detailedManifest() schema.ImageManifest {
	return schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-begin-test"),
		Labels:    systemLabels,
		App: &types.App{
			Exec:  types.Exec{"/bin/nethack4", "-D", "wizard"},
			User:  "0",
			Group: "0",
			Environment: types.Environment{
				types.EnvironmentVariable{
					Name:  "FOO",
					Value: "BAR",
				},
			},
			MountPoints: []types.MountPoint{
				types.MountPoint{
					Name:     *types.MustACName("nethack4-data"),
					Path:     "/root/nethack4-data",
					ReadOnly: true,
				},
			},
			Ports: []types.Port{
				types.Port{
					Name:     *types.MustACName("gopher"),
					Protocol: "tcp",
					Port:     70,
					Count:    1,
				},
			},
		},
		Annotations: types.Annotations{
			types.Annotation{
				Name:  *types.MustACIdentifier("author"),
				Value: "the acbuild devs",
			},
		},
		Dependencies: types.Dependencies{
			types.Dependency{
				ImageName: *types.MustACIdentifier("quay.io/gnu/hurd"),
			},
		},
	}
}

func runACBuildNoHist(workingDir string, args ...string) error {
	_, _, _, err := runACBuild(workingDir, append([]string{"--no-history"}, args...)...)
	return err
}

// runACBuild takes the workingDir and args to call ACBuild with, calls
// acbuild, and returns it's exit code, what it printed to stdout, what it
// printed to stderr, and an error in the event of a non-0 exit code.
func runACBuild(workingDir string, args ...string) (int, string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(acbuildBinPath, args...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = workingDir, &stdout, &stderr
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		code := exitErr.Sys().(syscall.WaitStatus).ExitStatus()
		acbuildOutput := fmt.Sprintf("stdout:\n%s\nstderr:\n%s\n", stdout.String(), stderr.String())
		return code, stdout.String(), stderr.String(), fmt.Errorf("non-zero exit code of %d: %s", code, acbuildOutput)
	}
	if err != nil {
		panic(err)
	}
	return 0, stdout.String(), stderr.String(), nil
}

func setUpTest(t *testing.T) string {
	tmpdir := mustTempDir()

	_, _, _, err := runACBuild(tmpdir, "begin")
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

func mustTempFile() *os.File {
	file, err := ioutil.TempFile("", "acbuild-test")
	if err != nil {
		panic(err)
	}
	return file
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
