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
	"path"
	"path/filepath"
	"syscall"

	"github.com/appc/spec/aci"
	"github.com/appc/spec/schema/types"

	"github.com/containers/build/util"
)

// Write will produce the resulting image from the current build context, saving
// it to the given path, optionally signing it.
func (a *ACBuild) Write(output string, overwrite bool) (err error) {
	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	if a.Mode == BuildModeAppC {
		man, err := util.GetManifest(a.CurrentImagePath)
		if err != nil {
			return err
		}

		if man.App != nil && len(man.App.Exec) == 0 {
			fmt.Fprintf(os.Stderr, "warning: exec command was never set.\n")
		}

		if man.Name == types.ACIdentifier(placeholdername) {
			return fmt.Errorf("can't write ACI, name was never set")
		}
	}

	fileFlags := os.O_CREATE | os.O_WRONLY

	_, err = os.Stat(output)
	switch {
	case os.IsNotExist(err):
		break
	case err != nil:
		return err
	default:
		if !overwrite {
			return fmt.Errorf("ACI already exists: %s", output)
		}
		fileFlags |= os.O_TRUNC
	}

	// open/create the image file
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

	// setup tar writer
	twriter := tar.NewWriter(gzwriter)
	defer twriter.Close()

	// create the aci writer
	switch a.Mode {
	case BuildModeAppC:
		man, err := util.GetManifest(a.CurrentImagePath)
		if err != nil {
			return err
		}
		aw := aci.NewImageWriter(*man, twriter)
		err = filepath.Walk(a.CurrentImagePath, aci.BuildWalker(a.CurrentImagePath, aw, nil))
		defer aw.Close()
		if err != nil {
			pathErr, ok := err.(*os.PathError)
			if !ok {
				fmt.Printf("not a path error!\n")
				return err
			}
			syscallErrno, ok := pathErr.Err.(syscall.Errno)
			if !ok {
				fmt.Printf("not a syscall errno!\n")
				return err
			}
			if pathErr.Op == "open" && syscallErrno != syscall.EACCES {
				return err
			}
			problemPath := pathErr.Path[len(path.Join(a.CurrentImagePath, aci.RootfsDir)):]
			return fmt.Errorf("%q: permission denied - call write as root", problemPath)
		}
	case BuildModeOCI:
		err = filepath.Walk(a.CurrentImagePath, util.PathWalker(twriter, a.CurrentImagePath))
		if err != nil {
			return err
		}
	}
	return nil
}
