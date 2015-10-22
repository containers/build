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
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

const (
	mountName = "html"
	mountPath = "/usr/html"

	mountName2 = "nethack4-data"
	mountPath2 = "/root/nethack4-data"
)

func manWithMounts(mounts []types.MountPoint) schema.ImageManifest {
	return schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:        nil,
			User:        "0",
			Group:       "0",
			MountPoints: mounts,
		},
		Labels: systemLabels,
	}
}

func TestAddMount(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	mounts := []types.MountPoint{
		types.MountPoint{
			Name: *types.MustACName(mountName),
			Path: mountPath,
		},
	}

	checkManifest(t, workingDir, manWithMounts(mounts))
	checkEmptyRootfs(t, workingDir)
}

func TestAddMountReadOnly(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath, "--read-only")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	mounts := []types.MountPoint{
		types.MountPoint{
			Name:     *types.MustACName(mountName),
			Path:     mountPath,
			ReadOnly: true,
		},
	}

	checkManifest(t, workingDir, manWithMounts(mounts))
	checkEmptyRootfs(t, workingDir)
}

func TestAdd2Mounts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath, "--read-only")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "mount", "add", mountName2, mountPath2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	mounts := []types.MountPoint{
		types.MountPoint{
			Name:     *types.MustACName(mountName),
			Path:     mountPath,
			ReadOnly: true,
		},
		types.MountPoint{
			Name: *types.MustACName(mountName2),
			Path: mountPath2,
		},
	}

	checkManifest(t, workingDir, manWithMounts(mounts))
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmMounts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath, "--read-only")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "mount", "add", mountName2, mountPath2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "mount", "rm", mountName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	mounts := []types.MountPoint{
		types.MountPoint{
			Name: *types.MustACName(mountName2),
			Path: mountPath2,
		},
	}

	checkManifest(t, workingDir, manWithMounts(mounts))
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmMount(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath, "--read-only")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "mount", "rm", mountName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}

func TestOverwriteMount(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "add", mountName, mountPath, "--read-only")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "mount", "add", mountName, mountPath2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	mounts := []types.MountPoint{
		types.MountPoint{
			Name: *types.MustACName(mountName),
			Path: mountPath2,
		},
	}

	checkManifest(t, workingDir, manWithMounts(mounts))
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentMount(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "mount", "rm", mountName)
	switch {
	case err == nil:
		t.Fatalf("mount remove didn't return an error when asked to remove nonexistent mount")
	case err.exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running mount remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}
