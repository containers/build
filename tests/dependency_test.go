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
	"strconv"
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

const (
	depName           = "example.com/app"
	depImageID        = "sha512-739d7ae77d9e"
	depLabel1Key      = "version"
	depLabel1Val      = "2"
	depLabel2Key      = "env"
	depLabel2Val      = "canary"
	depSize      uint = 9823749

	depName2 = "example.com/app2"
)

func newLabel(key, val string) string {
	return key + "=" + val
}

func manWithDeps(deps types.Dependencies) schema.ImageManifest {
	man := emptyManifest()
	man.Dependencies = deps
	return man
}

func TestAddDependency(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAddDependencyWithImageID(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName, "--image-id", depImageID)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	hash, err := types.NewHash(depImageID)
	if err != nil {
		panic(err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
			ImageID:   hash,
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAddDependencyWithLabels(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName,
		"--label", newLabel(depLabel1Key, depLabel1Val),
		"--label", newLabel(depLabel2Key, depLabel2Val))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
			Labels: types.Labels{
				types.Label{
					Name:  *types.MustACIdentifier(depLabel1Key),
					Value: depLabel1Val,
				},
				types.Label{
					Name:  *types.MustACIdentifier(depLabel2Key),
					Value: depLabel2Val,
				},
			},
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAddDependencyWithSize(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName, "--size", strconv.Itoa(int(depSize)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
			Size:      depSize,
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAdd2Dependencies(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName,
		"--image-id", depImageID,
		"--label", newLabel(depLabel1Key, depLabel1Val),
		"--label", newLabel(depLabel2Key, depLabel2Val),
		"--size", strconv.Itoa(int(depSize)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "dependency", "add", depName2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	hash, err := types.NewHash(depImageID)
	if err != nil {
		panic(err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
			ImageID:   hash,
			Labels: types.Labels{
				types.Label{
					Name:  *types.MustACIdentifier(depLabel1Key),
					Value: depLabel1Val,
				},
				types.Label{
					Name:  *types.MustACIdentifier(depLabel2Key),
					Value: depLabel2Val,
				},
			},
			Size: depSize,
		},
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName2),
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmDependencies(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName,
		"--image-id", depImageID,
		"--label", newLabel(depLabel1Key, depLabel1Val),
		"--label", newLabel(depLabel2Key, depLabel2Val),
		"--size", strconv.Itoa(int(depSize)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "dependency", "add", depName2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "dependency", "remove", depName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName2),
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmDependencie(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName,
		"--image-id", depImageID,
		"--label", newLabel(depLabel1Key, depLabel1Val),
		"--label", newLabel(depLabel2Key, depLabel2Val),
		"--size", strconv.Itoa(int(depSize)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "dependency", "remove", depName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifest())
	checkEmptyRootfs(t, workingDir)
}

func TestOverwriteDependency(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "dependency", "add", depName,
		"--image-id", depImageID,
		"--label", newLabel(depLabel1Key, depLabel1Val),
		"--label", newLabel(depLabel2Key, depLabel2Val),
		"--size", strconv.Itoa(int(depSize)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "dependency", "add", depName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	deps := types.Dependencies{
		types.Dependency{
			ImageName: *types.MustACIdentifier(depName),
		},
	}

	checkManifest(t, workingDir, manWithDeps(deps))
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentDependency(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	exitCode, _, _, err := runACBuild(workingDir, "--no-history", "dependency", "remove", depName)
	switch {
	case err == nil:
		t.Fatalf("dependency remove didn't return an error when asked to remove nonexistent dependency")
	case exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running dependency remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest())
	checkEmptyRootfs(t, workingDir)
}
