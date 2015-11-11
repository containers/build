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
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"

	"github.com/appc/acbuild/util"
)

// Write will produce the resulting ACI from the current build context, saving
// it to the given path, optionally signing it.
func (a *ACBuild) Write(output string, overwrite, sign bool, gpgflags []string) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	man, err := util.GetManifest(a.CurrentACIPath)
	if err != nil {
		return err
	}

	if man.App != nil && len(man.App.Exec) == 0 {
		fmt.Fprintf(os.Stderr, "warning: exec command was never set.\n")
	}

	if man.Name == types.ACIdentifier(placeholdername) {
		return fmt.Errorf("can't write ACI, name was never set")
	}

	fileFlags := os.O_CREATE | os.O_WRONLY

	ex, err := util.Exists(output)
	if err != nil {
		return err
	}
	if ex {
		if !overwrite {
			return fmt.Errorf("ACI already exists: %s", output)
		}
		fileFlags |= os.O_TRUNC
	}

	// open/create the aci file
	ofile, err := os.OpenFile(output, fileFlags, 0644)
	if err != nil {
		return err
	}
	defer ofile.Close()

	defer func() {
		// When write is done, if an error is encountered remove the partial
		// ACI that had been written.
		if err != nil {
			os.Remove(output)
			os.Remove(output + ".asc")
		}
	}()

	// setup compression
	gzwriter := gzip.NewWriter(ofile)
	defer gzwriter.Close()

	// create the aci writer
	aw := aci.NewImageWriter(*man, tar.NewWriter(gzwriter))
	err = filepath.Walk(a.CurrentACIPath, aci.BuildWalker(a.CurrentACIPath, aw, nil))
	defer aw.Close()
	if err != nil {
		return err
	}

	if sign {
		err = signACI(output, output+".asc", gpgflags)
		if err != nil {
			return err
		}
	}

	return nil
}

func signACI(acipath, signaturepath string, flags []string) error {
	if len(flags) == 0 {
		flags = []string{"--armor", "--yes"}
	}
	flags = append(flags, "--output", signaturepath, "--detach-sig", acipath)

	return util.Exec("gpg", flags...)
}
