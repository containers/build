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
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/appc/spec/aci"
	rkttar "github.com/coreos/rkt/pkg/tar"
	"github.com/coreos/rkt/pkg/user"
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

// ExtractImage will extract the contents of the image at path to the directory
// at dst. If fileMap is set, only files in it will be extracted.
func ExtractImage(path, dst string, fileMap map[string]struct{}) error {
	dst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dr, err := aci.NewCompressedReader(file)
	if err != nil {
		return fmt.Errorf("error decompressing image: %v", err)
	}
	defer dr.Close()

	uidRange := user.NewBlankUidRange()

	if os.Geteuid() == 0 {
		return rkttar.ExtractTar(dr, dst, true, uidRange, fileMap)
	}

	editor, err := rkttar.NewUidShiftingFilePermEditor(uidRange)
	if err != nil {
		return fmt.Errorf("error determining current user: %v", err)
	}
	return rkttar.ExtractTarInsecure(tar.NewReader(dr), dst, true, fileMap, editor)
}

func PathWalker(twriter *tar.Writer, tarSrcPath string) func(string, os.FileInfo, error) error {
	prefixLen := len(tarSrcPath + "/")
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path == tarSrcPath {
			return nil
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = path[prefixLen:]
		twriter.WriteHeader(hdr)

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		n, err := io.Copy(twriter, f)
		if err != nil {
			return err
		}
		if n != info.Size() {
			return fmt.Errorf("underwrite error")
		}
		return nil
	}
}
