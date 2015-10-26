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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/appc/acbuild/registry"
	"github.com/appc/acbuild/util"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/discovery"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

var (
	placeholdername = "acbuild-unnamed"
)

// Begin will start a new build, storing the untarred ACI the build operates on
// at a.CurrentACIPath. If start is the empty string, the build will begin with
// an empty ACI, otherwise the ACI stored at start will be used at the starting
// point.
func (a *ACBuild) Begin(start string, insecure bool) (err error) {
	ex, err := util.Exists(a.ContextPath)
	if err != nil {
		return err
	}
	if ex {
		return fmt.Errorf("build already in progress in this working dir")
	}

	err = os.MkdirAll(a.ContextPath, 0755)
	if err != nil {
		return err
	}

	if err = a.lock(); err != nil {
		return err
	}
	defer func() {
		if err1 := a.unlock(); err == nil {
			err = err1
		}
	}()

	if start != "" {
		err = os.MkdirAll(a.CurrentACIPath, 0755)
		if err != nil {
			return err
		}
		return a.beginFromImage(start, insecure)
	}

	err = os.MkdirAll(path.Join(a.CurrentACIPath, aci.RootfsDir), 0755)
	if err != nil {
		return err
	}

	acid, err := types.NewACIdentifier("acbuild-unnamed")
	if err != nil {
		return err
	}

	manifest := &schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *acid,
		App: &types.App{
			Exec:  nil,
			User:  "0",
			Group: "0",
		},
	}

	manblob, err := manifest.MarshalJSON()
	if err != nil {
		return err
	}

	manfile, err := os.Create(path.Join(a.CurrentACIPath, aci.ManifestFile))
	if err != nil {
		return err
	}

	_, err = manfile.Write(manblob)
	if err != nil {
		return err
	}

	err = manfile.Close()
	if err != nil {
		return err
	}

	return nil
}

func (a *ACBuild) beginFromImage(start string, insecure bool) error {
	// Check if we're starting with a file
	finfo, err := os.Stat(start)
	if err == nil {
		if finfo.IsDir() {
			return fmt.Errorf("provided starting ACI is a directory: %s", start)
		}
		return util.ExtractImage(start, a.CurrentACIPath, nil)
	} else if !os.IsNotExist(err) {
		return err
	}

	// Check if we're starting with a docker image
	if strings.HasPrefix(start, "docker://") {
		// TODO use docker2aci
		return fmt.Errorf("docker containers are currently unsupported")
	}

	// Perform meta discovery, download the ACI, and start with that.

	app, err := discovery.NewAppFromString(start)
	if err != nil {
		return err
	}
	labels, err := types.LabelsFromMap(app.Labels)
	if err != nil {
		return err
	}

	tmpDepStoreTarPath, err := ioutil.TempDir("", "acbuild-begin-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDepStoreTarPath)

	tmpDepStoreExpandedPath, err := ioutil.TempDir("", "acbuild-begin-expanded")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDepStoreExpandedPath)

	reg := registry.Registry{
		DepStoreTarPath:      tmpDepStoreTarPath,
		DepStoreExpandedPath: tmpDepStoreExpandedPath,
		Insecure:             insecure,
		Debug:                a.Debug,
	}

	err = reg.Fetch(app.Name, labels, 0, false)
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(tmpDepStoreTarPath)
	if err != nil {
		return err
	}

	if len(files) != 1 {
		var filelist string
		for _, file := range files {
			if filelist == "" {
				filelist = file.Name()
			} else {
				filelist = filelist + ", " + file.Name()
			}
		}
		panic("unexpected number of files in store after download: " + filelist)
	}

	return util.ExtractImage(path.Join(tmpDepStoreTarPath, files[0].Name()), a.CurrentACIPath, nil)
}
