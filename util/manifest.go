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

package util

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
)

// GetManifest will read in the manifest from the untarred ACI stored at acipath
func GetManifest(acipath string) (*schema.ImageManifest, error) {
	acifile, err := os.Open(path.Join(acipath, aci.ManifestFile))
	if err != nil {
		return nil, err
	}
	defer acifile.Close()

	manblob, err := ioutil.ReadAll(acifile)
	if err != nil {
		return nil, err
	}

	man := &schema.ImageManifest{}
	err = man.UnmarshalJSON(manblob)
	if err != nil {
		return nil, err
	}

	return man, nil
}

// ModifyManifest will read in the manifest from the untarred ACI stored at
// acipath, run the fn function (which is intended to modify the manifest), and
// then write the resulting manifest back to the file it was read from.
func ModifyManifest(fn func(*schema.ImageManifest), acipath string) error {
	man, err := GetManifest(acipath)
	if err != nil {
		return err
	}

	fn(man)

	blob, err := man.MarshalJSON()
	if err != nil {
		return err
	}

	manfile, err := os.OpenFile(path.Join(acipath, aci.ManifestFile),
		os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer manfile.Close()

	_, err = manfile.Write(blob)
	if err != nil {
		return err
	}

	return nil
}
