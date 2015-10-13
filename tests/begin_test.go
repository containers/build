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
	"io/ioutil"
	"path"
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/kylelemons/godebug/pretty"
)

func TestBegin(t *testing.T) {
	// Call begin
	tmpdir, err := setUpTest()
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer cleanUpTest(tmpdir)

	// Check that there are no files in the rootfs
	files, err := ioutil.ReadDir(path.Join(tmpdir, ".acbuild", "currentaci", aci.RootfsDir))
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(files) != 0 {
		t.Errorf("rootfs in aci contains files, should be empty")
	}

	// Check that the manifest is no bigger than it needs to be
	manblob, err := ioutil.ReadFile(path.Join(tmpdir, ".acbuild", "currentaci", aci.ManifestFile))
	if err != nil {
		t.Errorf("%v", err)
	}

	var man schema.ImageManifest

	err = man.UnmarshalJSON(manblob)
	if err != nil {
		t.Errorf("invalid manifest schema: %v", err)
	}

	expectedMan := schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:  nil,
			User:  "0",
			Group: "0",
		},
	}

	if str := pretty.Compare(man, expectedMan); str != "" {
		t.Errorf("unexpected manifest:\n%s", str)
	}
}
