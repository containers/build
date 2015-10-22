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
	annoName  = "authors"
	annoValue = "the acbuild developers"
)

var manWithOneAnno = schema.ImageManifest{
	ACKind:    schema.ImageManifestKind,
	ACVersion: schema.AppContainerVersion,
	Name:      *types.MustACIdentifier("acbuild-unnamed"),
	App: &types.App{
		Exec:  nil,
		User:  "0",
		Group: "0",
	},
	Annotations: types.Annotations{
		types.Annotation{
			Name:  *types.MustACIdentifier(annoName),
			Value: annoValue,
		},
	},
	Labels: systemLabels,
}

func TestAddAnnotation(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "annotation", "add", annoName, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, manWithOneAnno)
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmAnnotation(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "annotation", "add", annoName, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "annotation", "remove", annoName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}

func TestAddOverwriteAnnotation(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "annotation", "add", annoName, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "annotation", "add", annoName, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, manWithOneAnno)
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentAnnotation(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "annotation", "remove", annoName)
	switch {
	case err == nil:
		t.Fatalf("annotation remove didn't return an error when asked to remove nonexistent annotation")
	case err.exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running annotation remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmAnnotation(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	const suffix = "1"

	err := runACBuild(workingDir, "annotation", "add", annoName+suffix, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "annotation", "add", annoName, annoValue)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "annotation", "remove", annoName+suffix)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, manWithOneAnno)
	checkEmptyRootfs(t, workingDir)
}
