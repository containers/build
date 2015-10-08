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

// Write will produce the resulting ACI from the current build context, saving
// it to the given path, optionally signing it.
func Write(tmpaci, output string, overwrite, sign bool, gpgflags []string) error {
	man, err := util.GetManifest(tmpaci)
	if err != nil {
		return err
	}

	if man.App != nil && testEq(man.App.Exec, placeholderexec) {
		fmt.Fprintf(os.Stderr, "warning: exec command was never set.\n")
	}

	if man.Name == types.ACIdentifier(placeholdername) {
		return fmt.Errorf("can't write ACI, name was never set")
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

	err = filepath.Walk(tmpaci, aci.BuildWalker(tmpaci, aw, nil))

	aw.Close()
	ofile.Close()

	if err != nil {
		return err
	}

	if sign {
		err = signACI(output, output+".asc", gpgflags)
		if err != nil {
			os.Remove(output)
			os.Remove(output + ".asc")
			return err
		}
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

func signACI(acipath, signaturepath string, flags []string) error {
	if len(flags) == 0 {
		flags = []string{"--armor", "--yes"}
	}
	flags = append(flags, "--output", signaturepath, "--detach-sig", acipath)

	return util.Exec("gpg", flags...)
}
