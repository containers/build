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
	"path"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"

	"github.com/appc/acbuild/util"
)

// Copy will copy the directory/file at from to the path to inside the untarred
// ACI at acipath.
func Copy(acipath, from, to string) error {
	err := util.Exec("cp", "-r", from, path.Join(acipath, aci.RootfsDir, to))
	return err
}