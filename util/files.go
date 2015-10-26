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
	"fmt"
	"os"
	"path/filepath"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/tar"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/uid"
)

// RmAndMkdir will remove anything at path if it exists, and then create a
// directory at path.
func RmAndMkdir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// Exists will return whether or not anything exists at path
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// UnTar will extract the contents at the tar file at tarpath to the directory
// at dst. If fileMap is set, only files in it will be extracted.
func UnTar(tarpath, dst string, fileMap map[string]struct{}) error {
	dst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	tarfile, err := os.Open(tarpath)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	dr, err := aci.NewCompressedReader(tarfile)
	if err != nil {
		return fmt.Errorf("error decompressing image: %v", err)
	}
	defer dr.Close()

	return tar.ExtractTar(dr, dst, true, uid.NewBlankUidRange(), fileMap)
}
