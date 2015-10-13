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
	"archive/tar"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"

	"github.com/appc/acbuild/util"
)

// End will end the build, storing the untarred ACI stored at tmpaci into
// output. The files at contextpath will be removed, to end the build. If
// overwrite is true, an error will not be thrown if output already exists.
func End(tmpaci, output, contextpath string, overwrite bool) error {
	man, err := util.GetManifest(tmpaci)
	if err != nil {
		return err
	}

	if man.App != nil && testEq(man.App.Exec, placeholderexec) {
		return fmt.Errorf("can't end build, exec command was never set")
	}

	if man.Name == types.ACIdentifier(placeholdername) {
		return fmt.Errorf("can't end build, name was never set")
	}

	ex, err := util.Exists(output)
	if err != nil {
		return err
	}
	if ex {
		if !overwrite {
			return fmt.Errorf("ACI already exists: %s", output)
		}
		err := os.Remove(output)
		if err != nil {
			return err
		}
	}

	ofile, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	aw := aci.NewImageWriter(*man, tar.NewWriter(ofile))
	defer aw.Close()

	err = filepath.Walk(tmpaci, aci.BuildWalker(tmpaci, aw, func(hdr *tar.Header) bool { return true }))
	if err != nil {
		return err
	}

	err = os.RemoveAll(contextpath)
	if err != nil {
		return err
	}

	return nil
}

func testEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, elem := range a {
		if elem != b[i] {
			return false
		}
	}
	return true
}
