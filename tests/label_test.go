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
	labelName = "version"
	labelVal  = "1.0.0"

	labelName2 = "foo"
	labelVal2  = "bar"
)

func manWithLabels(labels types.Labels) schema.ImageManifest {
	man := emptyManifest()
	man.Labels = append(man.Labels, labels...)
	return man
}

func TestAddLabel(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "label", "add", labelName, labelVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	labels := types.Labels{
		types.Label{
			Name:  *types.MustACIdentifier(labelName),
			Value: labelVal,
		},
	}

	checkManifest(t, workingDir, manWithLabels(labels))
	checkEmptyRootfs(t, workingDir)
}

func TestAddTwoLabels(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "label", "add", labelName, labelVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "label", "add", labelName2, labelVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	labels := types.Labels{
		types.Label{
			Name:  *types.MustACIdentifier(labelName),
			Value: labelVal,
		},
		types.Label{
			Name:  *types.MustACIdentifier(labelName2),
			Value: labelVal2,
		},
	}

	checkManifest(t, workingDir, manWithLabels(labels))
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmLabels(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "label", "add", labelName, labelVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "label", "add", labelName2, labelVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "label", "rm", labelName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	labels := types.Labels{
		types.Label{
			Name:  *types.MustACIdentifier(labelName2),
			Value: labelVal2,
		},
	}

	checkManifest(t, workingDir, manWithLabels(labels))
	checkEmptyRootfs(t, workingDir)
}

func TestAddOverwriteLabel(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "label", "add", labelName, labelVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "label", "add", labelName, labelVal2)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	labels := types.Labels{
		types.Label{
			Name:  *types.MustACIdentifier(labelName),
			Value: labelVal2,
		},
	}

	checkManifest(t, workingDir, manWithLabels(labels))
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmLabel(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuildNoHist(workingDir, "label", "add", labelName, labelVal)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuildNoHist(workingDir, "label", "rm", labelName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifest())
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentLabel(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	exitCode, _, _, err := runACBuild(workingDir, "--no-history", "label", "rm", labelName)
	switch {
	case err == nil:
		t.Fatalf("label remove didn't return an error when asked to remove nonexistent label")
	case exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running label remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest())
	checkEmptyRootfs(t, workingDir)
}
