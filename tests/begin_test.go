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
	"archive/tar"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestBeginEmpty(t *testing.T) {
	tmpdir := mustTempDir()
	defer cleanUpTest(tmpdir)

	err := runAcbuild(tmpdir, "begin")
	if err != nil {
		t.Fatalf("%v", err)
	}

	wim := schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:  nil,
			User:  "0",
			Group: "0",
		},
	}

	checkMinimalContainer(t, path.Join(tmpdir, ".acbuild", "currentaci"), wim)
}

func TestBeginLocalACI(t *testing.T) {
	wim := schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-begin-test"),
		Labels: types.Labels{
			types.Label{
				Name:  *types.MustACIdentifier("version"),
				Value: "9001",
			},
		},
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

	manblob, err := wim.MarshalJSON()
	if err != nil {
		panic(err)
	}

	tmpexpandedaci := mustTempDir()
	defer os.RemoveAll(tmpexpandedaci)

	err = ioutil.WriteFile(path.Join(tmpexpandedaci, aci.ManifestFile), manblob, 0644)
	if err != nil {
		panic(err)
	}

	err = os.Mkdir(path.Join(tmpexpandedaci, aci.RootfsDir), 0755)
	if err != nil {
		panic(err)
	}

	tmpaci, err := ioutil.TempFile("", "acbuild-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpaci.Name())

	aw := aci.NewImageWriter(wim, tar.NewWriter(tmpaci))
	err = filepath.Walk(tmpexpandedaci, aci.BuildWalker(tmpexpandedaci, aw, nil))
	aw.Close()
	if err != nil {
		panic(err)
	}
	tmpaci.Close()

	tmpdir := mustTempDir()
	defer cleanUpTest(tmpdir)

	err = runAcbuild(tmpdir, "begin", tmpaci.Name())
	if err != nil {
		t.Fatalf("%v", err)
	}

	checkMinimalContainer(t, path.Join(tmpdir, ".acbuild", "currentaci"), wim)
}

func checkMinimalContainer(t *testing.T, acipath string, expectedManifest schema.ImageManifest) {
	// Check that there are no files in the rootfs
	files, err := ioutil.ReadDir(path.Join(acipath, aci.RootfsDir))
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(files) != 0 {
		t.Errorf("rootfs in aci contains files, should be empty")
	}

	// Check that the manifest is no bigger than it needs to be
	manblob, err := ioutil.ReadFile(path.Join(acipath, aci.ManifestFile))
	if err != nil {
		t.Errorf("%v", err)
	}

	var man schema.ImageManifest

	err = man.UnmarshalJSON(manblob)
	if err != nil {
		t.Errorf("invalid manifest schema: %v", err)
	}

	if str := pretty.Compare(man, expectedManifest); str != "" {
		t.Errorf("unexpected manifest:\n%s", str)
	}
}
