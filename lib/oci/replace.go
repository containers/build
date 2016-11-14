// Copyright 2016 The acbuild Authors
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

package oci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Replace will replace the config in the expanded OCI image stored at
// a.CurrentImagePath with the new config stored at configPath
func (i *Image) Replace(configPath string) error {
	finfo, err := os.Stat(configPath)
	switch {
	case os.IsNotExist(err):
		return fmt.Errorf("no such file or directory: %s", configPath)
	case err != nil:
		return err
	case finfo.IsDir():
		return fmt.Errorf("%s is a directory", configPath)
	default:
		break
	}

	confblob, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(confblob, &i.config)
	if err != nil {
		return err
	}

	return i.save()
}
