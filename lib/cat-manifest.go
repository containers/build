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
	"io"
	"io/ioutil"
	"os"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"

	"github.com/appc/acbuild/util"
)

// CatManifest will print to stdout the manifest from the expanded ACI stored
// at a.CurrentACIPath, optionally inserting whitespace to make it more human
// readable.
func (a *ACBuild) CatManifest(prettyPrint bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	man, err := util.GetManifest(a.CurrentACIPath)
	if err != nil {
		return err
	}

	return util.PrintManifest(man, prettyPrint)
}

// CatManifest will print to stdout the manifest from the ACI stored at
// aciPath, optionally inserting whitespace to make it more human readable.
func CatManifest(aciPath string, prettyPrint bool) (err error) {
	finfo, err := os.Stat(aciPath)
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("no such file or directory: %s", aciPath)
	case err != nil:
		return err
	case finfo.IsDir():
		return fmt.Errorf("%s is a directory, not an ACI", aciPath)
	default:
		break
	}

	file, err := os.Open(aciPath)
	if err != nil {
		return err
	}
	defer file.Close()

	tr, err := aci.NewCompressedTarReader(file)
	if err != nil {
		return fmt.Errorf("error decompressing image: %v", err)
	}
	defer tr.Close()

	for {
		hdr, err := tr.Next()
		switch {
		case err == io.EOF:
			return fmt.Errorf("manifest not found in ACI %s", aciPath)
		case err != nil:
			return err
		case hdr.Name == "manifest":
			manblob, err := ioutil.ReadAll(tr)
			if err != nil {
				return err
			}
			var man schema.ImageManifest
			err = man.UnmarshalJSON(manblob)
			if err != nil {
				return err
			}
			return util.PrintManifest(&man, prettyPrint)
		}
	}
}
