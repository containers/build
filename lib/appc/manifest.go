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

package appc

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
)

var (
	// ErrNotFound is returned when acbuild is asked to remove an element from a
	// list and the element is not present in the list
	ErrNotFound = fmt.Errorf("element to be removed does not exist in this ACI")
)

// Manifest is a struct with an open handle to a manifest that it can manipulate
type Manifest struct {
	aciPath  string
	manifest *schema.ImageManifest
}

// LoadManifest will read in the manifest from an untarred ACI on disk at
// location aciPath, and return a new Manifest struct to manipulate it.
func LoadManifest(aciPath string) (*Manifest, error) {
	manFile, err := os.OpenFile(path.Join(aciPath, aci.ManifestFile), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer manFile.Close()

	manblob, err := ioutil.ReadAll(manFile)
	if err != nil {
		return nil, err
	}

	man := &schema.ImageManifest{}
	err = man.UnmarshalJSON(manblob)
	if err != nil {
		return nil, err
	}

	return &Manifest{aciPath, man}, nil
}

// save commits changes to m's manifest to disk
func (m *Manifest) save() error {
	blob, err := m.manifest.MarshalJSON()
	if err != nil {
		return err
	}

	manFile, err := os.OpenFile(path.Join(m.aciPath, aci.ManifestFile), os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer manFile.Close()

	err = manFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = manFile.Write(blob)
	if err != nil {
		return err
	}

	return nil
}

// Get returns the manifest currently being manipulated
func (m *Manifest) Get() *schema.ImageManifest {
	return m.manifest
}

// Print will print out the current manifest to stdout, optionally inserting
// whitespace to improve readability
func (m *Manifest) Print(w io.Writer, prettyPrint, printConfig bool) error {
	if printConfig {
		return fmt.Errorf("can't print config, appc has no image configs")
	}
	var manblob []byte
	var err error
	if prettyPrint {
		manblob, err = json.MarshalIndent(m.manifest, "", "    ")
	} else {
		manblob, err = m.manifest.MarshalJSON()
	}
	if err != nil {
		return err
	}
	manblob = append(manblob, '\n')
	n, err := w.Write(manblob)
	if err != nil {
		return err
	}
	if n < len(manblob) {
		return fmt.Errorf("short write")
	}
	return nil
}

// newManifestApp will generate a valid minimal types.App for use in a
// schema.ImageManifest. This is necessary as placing a completely empty
// types.App into a manifest will result in an invalid manifest.
func newManifestApp() *types.App {
	return &types.App{
		User:  "0",
		Group: "0",
	}
}
