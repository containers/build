// Copyright 2016 The appc Authors
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
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

const (
	isolatorName  = "os/linux/capabilities-retain-set"
	isolatorValue = `{ "set": ["CAP_NET_BIND_SERVICE"] }`
)

func manWithIsolators(isolators types.Isolators) schema.ImageManifest {
	man := emptyManifestWithApp()
	man.App.Isolators = isolators
	return man
}

func TestAddIsolator(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	isoFile, err := ioutil.TempFile("", "acbuild-test")
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	_, err = isoFile.Write([]byte(isolatorValue))
	isoFile.Close()
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer os.RemoveAll(isoFile.Name())

	err = runACBuildNoHist(workingDir, "isolator", "add", isolatorName, isoFile.Name())
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	valueBlob := json.RawMessage(isolatorValue)
	i := &types.Isolator{
		Name:     *types.MustACIdentifier(isolatorName),
		ValueRaw: &valueBlob,
	}
	blob, err := json.Marshal(i)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	err = i.UnmarshalJSON(blob)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, manWithIsolators(types.Isolators{*i}))
	checkEmptyRootfs(t, workingDir)
}
