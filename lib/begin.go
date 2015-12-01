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
	"net"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/appc/acbuild/registry"
	"github.com/appc/acbuild/util"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/aci"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/discovery"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/fileutil"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/coreos/rkt/pkg/uid"
)

var (
	placeholdername = "acbuild-unnamed"
)

// Begin will start a new build, storing the untarred ACI the build operates on
// at a.CurrentACIPath. If start is the empty string, the build will begin with
// an empty ACI, otherwise the ACI stored at start will be used at the starting
// point.
func (a *ACBuild) Begin(start string, insecure bool) (err error) {
	_, err = os.Stat(a.ContextPath)
	switch {
	case os.IsNotExist(err):
		break
	case err != nil:
		return err
	default:
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

	defer func() {
		// If there was an error while beginning, we don't want to produce an
		// unexpected build context
		if err != nil {
			os.RemoveAll(a.ContextPath)
		}
	}()

	if start != "" {
		err = os.MkdirAll(a.CurrentACIPath, 0755)
		if err != nil {
			return err
		}
		if start[0] == '.' || start[0] == '/' {
			finfo, err := os.Stat(start)
			switch {
			case os.IsNotExist(err):
				return fmt.Errorf("no such file or directory: %s", start)
			case err != nil:
				return err
			case finfo.IsDir():
				a.beginFromLocalDirectory(start)
			default:
				return a.beginFromLocalImage(start)
			}
		} else {
			return a.beginFromRemoteImage(start, insecure)
		}
	}
	return a.beginWithEmptyACI()
}

func (a *ACBuild) beginFromLocalImage(start string) error {
	finfo, err := os.Stat(start)
	if err != nil {
		return err
	}
	if finfo.IsDir() {
		return fmt.Errorf("provided starting ACI is a directory: %s", start)
	}
	return util.ExtractImage(start, a.CurrentACIPath, nil)
}

func (a *ACBuild) beginFromLocalDirectory(start string) error {
	err := os.MkdirAll(a.CurrentACIPath, 0755)
	if err != nil {
		return err
	}

	err = fileutil.CopyTree(start, path.Join(a.CurrentACIPath, aci.RootfsDir), uid.NewBlankUidRange())
	if err != nil {
		return err
	}

	return a.writeEmptyManifest()
}

func (a *ACBuild) beginWithEmptyACI() error {
	err := os.MkdirAll(path.Join(a.CurrentACIPath, aci.RootfsDir), 0755)
	if err != nil {
		return err
	}

	return a.writeEmptyManifest()
}

func (a *ACBuild) writeEmptyManifest() error {
	acid, err := types.NewACIdentifier("acbuild-unnamed")
	if err != nil {
		return err
	}

	archlabel, err := types.NewACIdentifier("arch")
	if err != nil {
		return err
	}

	oslabel, err := types.NewACIdentifier("os")
	if err != nil {
		return err
	}

	manifest := &schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *acid,
		Labels: types.Labels{
			types.Label{
				*archlabel,
				runtime.GOARCH,
			},
			types.Label{
				*oslabel,
				runtime.GOOS,
			},
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

func (a *ACBuild) beginFromRemoteImage(start string, insecure bool) error {
	// Check if we're starting with a docker image
	if strings.HasPrefix(start, "docker://") {
		// TODO use docker2aci
		return fmt.Errorf("docker containers are currently unsupported")
	}

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
		if urlerr, ok := err.(*url.Error); ok {
			if operr, ok := urlerr.Err.(*net.OpError); ok {
				if dnserr, ok := operr.Err.(*net.DNSError); ok {
					if dnserr.Err == "no such host" {
						return fmt.Errorf("unknown host when fetching image, check your connection and local file paths must start with '/' or '.'")
					}
				}
			}
		}
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
