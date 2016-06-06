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

	target := path.Join(a.CurrentACIPath, aci.RootfsDir, to)

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

	target := path.Join(a.CurrentACIPath, aci.RootfsDir, to)

	dir, _ := path.Split(target)
	if dir != "" {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return fileutil.CopyTree(from, target, user.NewBlankUidRange())
}
