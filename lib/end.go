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

// End will stop the current build. An error will be returned if no build is in
// progress.
func (a *ACBuild) End() error {
	ok, err := util.Exists(a.ContextPath)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("build not in progress")
	}

	if err = a.lock(); err != nil {
		return err
	}

	err = os.RemoveAll(a.ContextPath)
	if err != nil {
		return err
	}

	return nil
}
