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
	"os"
	"path"

	"github.com/appc/spec/aci"
	"github.com/coreos/rkt/pkg/fileutil"
	"github.com/coreos/rkt/pkg/user"

	"github.com/containers/build/lib/oci"
	"github.com/containers/build/util"
)

// CopyToDir will copy all elements specified in the froms slice into the
// directory inside the current ACI specified by the to string.
func (a *ACBuild) CopyToDir(froms []string, to string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	switch a.Mode {
	case BuildModeAppC:
		return a.copyToDirAppC(froms, to)
	case BuildModeOCI:
		return a.copyToDirOCI(froms, to)
	}
	return fmt.Errorf("unknown build mode: %s", a.Mode)
}

func (a *ACBuild) copyToDirAppC(froms []string, to string) error {
	target := path.Join(a.CurrentImagePath, aci.RootfsDir, to)

	targetInfo, err := os.Stat(target)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(target, 0755)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	case !targetInfo.IsDir():
		return fmt.Errorf("target %q is not a directory", to)
	}

	for _, from := range froms {
		_, file := path.Split(from)
		tmptarget := path.Join(target, file)
		err := fileutil.CopyTree(from, tmptarget, user.NewBlankUidRange())
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *ACBuild) expandTopOCILayer() (string, error) {
	var topLayerID string
	switch ociMan := a.man.(type) {
	case *oci.Image:
		layerIDs := ociMan.GetLayerDigests()
		if len(layerIDs) > 0 {
			topLayerID = layerIDs[len(layerIDs)-1]
		}
	default:
		return "", fmt.Errorf("internal error: mismatched manifest type and build mode???")
	}

	var targetPath string
	if topLayerID == "" {
		var err error
		targetPath, err = util.OCINewExpandedLayer(a.OCIExpandedBlobsPath)
		if err != nil {
			return "", err
		}
	} else {
		topLayerAlgo, topLayerHash, err := util.SplitOCILayerID(topLayerID)
		if err != nil {
			return "", err
		}

		err = util.OCIExtractLayers([]string{topLayerID}, a.CurrentImagePath, a.OCIExpandedBlobsPath)
		if err != nil {
			return "", err
		}
		targetPath = path.Join(a.OCIExpandedBlobsPath, topLayerAlgo, topLayerHash)
	}
	return targetPath, nil
}

func (a *ACBuild) copyToDirOCI(froms []string, to string) error {
	currentLayer, err := a.expandTopOCILayer()
	if err != nil {
		return err
	}
	targetPath := path.Join(currentLayer, to)

	targetInfo, err := os.Stat(targetPath)
	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			return err
		}
	case err != nil:
		return err
	case !targetInfo.IsDir():
		return fmt.Errorf("target %q is not a directory", to)
	}

	for _, from := range froms {
		_, file := path.Split(from)
		tmptarget := path.Join(targetPath, file)
		err := fileutil.CopyTree(from, tmptarget, user.NewBlankUidRange())
		if err != nil {
			return err
		}
	}

	return a.rehashAndStoreOCIBlob(currentLayer, false)
}

// CopyToTarget will copy a single file/directory from the from string to the
// path specified by the to string inside the current ACI.
func (a *ACBuild) CopyToTarget(from string, to string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	switch a.Mode {
	case BuildModeAppC:
		return a.copyToTargetAppC(from, to)
	case BuildModeOCI:
		return a.copyToTargetOCI(from, to)
	}
	return fmt.Errorf("unknown build mode: %s", a.Mode)
}

func (a *ACBuild) copyToTargetAppC(from string, to string) error {
	target := path.Join(a.CurrentImagePath, aci.RootfsDir, to)

	dir, _ := path.Split(target)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return fileutil.CopyTree(from, target, user.NewBlankUidRange())
}

func (a *ACBuild) copyToTargetOCI(from string, to string) error {
	targetPath, err := a.expandTopOCILayer()
	if err != nil {
		return err
	}
	target := path.Join(targetPath, to)

	dir, _ := path.Split(target)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	err = fileutil.CopyTree(from, target, user.NewBlankUidRange())
	if err != nil {
		return err
	}

	return a.rehashAndStoreOCIBlob(targetPath, false)

}
