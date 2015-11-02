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
	envName = "FOO"
	envVal  = "BAR"

	envName2 = "BOO"
	envVal2  = "FAR"
)

func manWithEnv(env types.Environment) schema.ImageManifest {
	return schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:        nil,
			User:        "0",
			Group:       "0",
			Environment: env,
		},
		Labels: systemLabels,
	}
}

func TestAddEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "add", envName, envVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	env := types.Environment{
		types.EnvironmentVariable{
			Name:  envName,
			Value: envVal,
		},
	}

	checkManifest(t, workingDir, manWithEnv(env))
	checkEmptyRootfs(t, workingDir)
}

func TestAddTwoEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "add", envName, envVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "environment", "add", envName2, envVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	env := types.Environment{
		types.EnvironmentVariable{
			Name:  envName,
			Value: envVal,
		},
		types.EnvironmentVariable{
			Name:  envName2,
			Value: envVal2,
		},
	}

	checkManifest(t, workingDir, manWithEnv(env))
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "add", envName, envVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "environment", "add", envName2, envVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "environment", "rm", envName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	env := types.Environment{
		types.EnvironmentVariable{
			Name:  envName2,
			Value: envVal2,
		},
	}

	checkManifest(t, workingDir, manWithEnv(env))
	checkEmptyRootfs(t, workingDir)
}

func TestAddOverwriteEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "add", envName, envVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "environment", "add", envName, envVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	env := types.Environment{
		types.EnvironmentVariable{
			Name:  envName,
			Value: envVal2,
		},
	}

	checkManifest(t, workingDir, manWithEnv(env))
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "add", envName, envVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "environment", "rm", envName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifestWithApp)
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentEnv(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "environment", "rm", envName)
	switch {
	case err == nil:
		t.Fatalf("environment remove didn't return an error when asked to remove nonexistent environment variable")
	case err.exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running environment remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}
