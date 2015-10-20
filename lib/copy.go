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
	"os"
	"path"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/fileutil"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/uid"
)

// Copy will copy the directory/file at from to the path to inside the untarred
// ACI at a.CurrentACIPath.
func (a *ACBuild) Copy(from, to string) (err error) {
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

	return fileutil.CopyTree(from, target, uid.NewBlankUidRange())
}
