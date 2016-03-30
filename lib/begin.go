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
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"github.com/appc/acbuild/registry"
	"github.com/appc/acbuild/util"

	docker2aci "github.com/appc/docker2aci/lib"
	"github.com/appc/docker2aci/lib/common"
	"github.com/appc/spec/aci"
	"github.com/appc/spec/discovery"
	"github.com/appc/spec/schema"
	"github.com/appc/spec/schema/types"
	"github.com/coreos/rkt/pkg/fileutil"
	"github.com/coreos/rkt/pkg/user"
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
			dockerPrefix := "docker://"
			if strings.HasPrefix(start, dockerPrefix) {
				start = strings.TrimPrefix(start, dockerPrefix)
				return a.beginFromRemoteDockerImage(start, insecure)
			}
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
	err = util.ExtractImage(start, a.CurrentACIPath, nil)
	if err != nil {
		return err
	}

	for _, file := range []struct {
		FileName string
		FilePath string
	}{
		{"manifest file", path.Join(a.CurrentACIPath, aci.ManifestFile)},
		{"rootfs directory", path.Join(a.CurrentACIPath, aci.RootfsDir)},
	} {
		_, err = os.Stat(file.FilePath)
		switch {
		case os.IsNotExist(err):
			fmt.Fprintf(os.Stderr, "%s is missing, assuming build is beginning with a tar of a rootfs\n", file.FileName)
			return a.startedFromTar()
		case err != nil:
			return err
		}
	}
	return nil
}

func (a *ACBuild) startedFromTar() error {
	tmpPath := path.Join(a.ContextPath, "rootfs")
	err := os.Rename(a.CurrentACIPath, tmpPath)
	if err != nil {
		return err
	}
	err = a.beginWithEmptyACI()
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(a.CurrentACIPath, aci.RootfsDir))
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, path.Join(a.CurrentACIPath, aci.RootfsDir))
}

func (a *ACBuild) beginFromLocalDirectory(start string) error {
	err := os.MkdirAll(a.CurrentACIPath, 0755)
	if err != nil {
		return err
	}

	err = fileutil.CopyTree(start, path.Join(a.CurrentACIPath, aci.RootfsDir), user.NewBlankUidRange())
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

	archvalue := runtime.GOARCH
	if runtime.GOOS == "linux" && (archvalue == "arm" || archvalue == "arm64") {
		var x uint32 = 0x01020304
		test := *(*byte)(unsafe.Pointer(&x))
		switch {
		case test == 0x01 && archvalue == "arm":
			archvalue = "armv7b"
		case test == 0x04 && archvalue == "arm":
			archvalue = "armv7l"
		case test == 0x01 && archvalue == "arm64":
			archvalue = "aarch64_be"
		case test == 0x04 && archvalue == "arm64":
			archvalue = "aarch64"
		}
	}

	manifest := &schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *acid,
		Labels: types.Labels{
			types.Label{
				*archlabel,
				archvalue,
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

func (a *ACBuild) beginFromRemoteDockerImage(start string, insecure bool) (err error) {
	outputDir, err := ioutil.TempDir("", "acbuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(outputDir)

	tempDir, err := ioutil.TempDir("", "acbuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	insecureConf := common.InsecureConfig{
		SkipVerify: insecure,
		AllowHTTP:  insecure,
	}

	config := docker2aci.RemoteConfig{
		CommonConfig: docker2aci.CommonConfig{
			Squash:      true,
			OutputDir:   outputDir,
			TmpDir:      tempDir,
			Compression: common.GzipCompression,
		},
		Username: "",
		Password: "",
		Insecure: insecureConf,
	}
	renderedACIs, err := docker2aci.ConvertRemoteRepo(start, config)
	if err != nil {
		return err
	}
	if len(renderedACIs) > 1 {
		return fmt.Errorf("internal error: docker2aci didn't squash the image")
	}
	if len(renderedACIs) == 0 {
		return fmt.Errorf("internal error: docker2aci didn't produce any images")
	}
	absRenderedACI, err := filepath.Abs(renderedACIs[0])
	if err != nil {
		return err
	}

	return util.ExtractImage(absRenderedACI, a.CurrentACIPath, nil)
}
