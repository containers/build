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

package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/schema"
)

// ReplaceManifest will replace the manifest in the expanded ACI stored at
// a.CurrentACIPath with the new manifest stored at manifestPath
func (a *ACBuild) ReplaceManifest(manifestPath string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	finfo, err := os.Stat(manifestPath)
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("no such file or directory: %s", manifestPath)
	case err != nil:
		return err
	case finfo.IsDir():
		return fmt.Errorf("%s is a directory", manifestPath)
	default:
		break
	}

	manblob, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	// Marshal and Unmarshal the manifest to assert that it's valid and to
	// strip any whitespace

	var man schema.ImageManifest
	err = man.UnmarshalJSON(manblob)
	if err != nil {
		return err
	}

	manblob, err = man.MarshalJSON()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(a.CurrentACIPath, aci.ManifestFile), manblob, 0755)
}
