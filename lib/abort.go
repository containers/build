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

	"github.com/appc/acbuild/util"
)

// Abort will abort the current build, given the path that the build resources
// are stored at. An error will be returned if no build is in progress.
func Abort(contextpath string) error {
	ok, err := util.Exists(contextpath)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("build not in progress")
	}

	err = os.RemoveAll(contextpath)
	if err != nil {
		return err
	}

	return nil
}
