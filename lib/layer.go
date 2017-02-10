// Copyright 2017 The acbuild Authors
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

	"github.com/containers/build/util"
)

func (a *ACBuild) NewLayer() (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	if a.Mode != BuildModeOCI {
		return fmt.Errorf("adding layers is currently only supported in OCI builds")
	}

	newLayer, err := util.OCINewExpandedLayer(a.OCIExpandedBlobsPath)
	if err != nil {
		return err
	}
	return a.rehashAndStoreOCIBlob(newLayer, true)
}
