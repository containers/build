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
	"fmt"
	"io/ioutil"
	"os"
)

// Replace will replace the manifest in the expanded ACI stored at
// a.CurrentImagePath with the new manifest stored at manifestPath
func (m *Manifest) Replace(manifestPath string) error {
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

	err = m.manifest.UnmarshalJSON(manblob)
	if err != nil {
		return err
	}

	return m.save()
}
