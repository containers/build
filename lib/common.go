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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/containers/build/lib/appc"
	"github.com/containers/build/lib/oci"
	"github.com/containers/build/util"
)

const OCISchemaVersion = 2

const defaultWorkPath = ".acbuild"

type OCILayout struct {
	imageLayoutVersion string `json:"imageLayoutVersion"`
}

var OCILayoutValue = OCILayout{"1.0.0"}

// BuildMode represents which image spec is being followed during a build, AppC
// or OCI
type BuildMode string

const (
	BuildModeAppC = BuildMode("appc")
	BuildModeOCI  = BuildMode("oci")
)

var (
	errNoBuildInProgress = fmt.Errorf("no build in progress in this working dir - try \"acbuild begin\"")
)

// ACBuild contains all the information for a current build. Once an ACBuild
// has been created, the functions available on it will perform different
// actions in the build, like updating a dependency or writing a finished ACI.
type ACBuild struct {
	ContextPath          string
	LockPath             string
	CurrentImagePath     string
	DepStoreTarPath      string
	DepStoreExpandedPath string
	OverlayTargetPath    string
	OverlayWorkPath      string
	BuildModePath        string
	OCIExpandedBlobsPath string
	Debug                bool
	Mode                 BuildMode

	man      Manifest
	lockFile *os.File
}

// NewACBuild returns a new ACBuild struct with sane defaults for all of the
// different paths
func NewACBuild(cwd string, debug bool, buildMode BuildMode) (*ACBuild, error) {
	a := &ACBuild{
		ContextPath:          path.Join(cwd, defaultWorkPath),
		LockPath:             path.Join(cwd, defaultWorkPath, "lock"),
		CurrentImagePath:     path.Join(cwd, defaultWorkPath, "currentaci"),
		DepStoreTarPath:      path.Join(cwd, defaultWorkPath, "depstore-tar"),
		DepStoreExpandedPath: path.Join(cwd, defaultWorkPath, "depstore-expanded"),
		OverlayTargetPath:    path.Join(cwd, defaultWorkPath, "target"),
		OverlayWorkPath:      path.Join(cwd, defaultWorkPath, "work"),
		BuildModePath:        path.Join(cwd, defaultWorkPath, "buildMode"),
		OCIExpandedBlobsPath: path.Join(cwd, defaultWorkPath, "ociblobs"),
		Debug:                debug,
		Mode:                 buildMode,
	}
	// This might fail, and that's ok (maybe the build hasn't started yet)
	a.loadManifest()
	return a, nil
}

func (a *ACBuild) loadManifest() error {
	var err error
	switch a.Mode {
	case BuildModeAppC:
		a.man, err = appc.LoadManifest(a.CurrentImagePath)
	case BuildModeOCI:
		a.man, err = oci.LoadImage(a.CurrentImagePath)
	}
	if err != nil {
		_, serr := os.Stat(a.ContextPath)
		if !os.IsNotExist(serr) {
			// If the context path exists, then this build was started and we
			// shouldn't have failed. Let's error.
			return fmt.Errorf("error loading manifest: %v\n", err)
		}
	}
	return nil
}

func (a *ACBuild) lock() error {
	_, err := os.Stat(a.ContextPath)
	switch {
	case os.IsNotExist(err):
		return errNoBuildInProgress
	case err != nil:
		return err
	}

	if a.lockFile != nil {
		return fmt.Errorf("lock already held by this ACBuild")
	}

	a.lockFile, err = os.OpenFile(a.LockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		switch err1 := err.(type) {
		case *os.PathError:
			switch err2 := err1.Err.(type) {
			case syscall.Errno:
				if err2 == syscall.EACCES {
					err = fmt.Errorf("permission denied: please run this as a user with appropriate privileges\n")
				}
			}
		}
		return err
	}

	err = syscall.Flock(int(a.lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		if err == syscall.EWOULDBLOCK {
			return fmt.Errorf("lock already held - is another acbuild running in this working dir?")
		}
		return err
	}

	return nil
}

func (a *ACBuild) unlock() error {
	if a.lockFile == nil {
		return fmt.Errorf("lock isn't held by this ACBuild")
	}

	err := syscall.Flock(int(a.lockFile.Fd()), syscall.LOCK_UN)
	if err != nil {
		return err
	}

	err = a.lockFile.Close()
	if err != nil {
		return err
	}
	a.lockFile = nil

	err = os.Remove(a.LockPath)
	if err != nil {
		return err
	}

	return nil
}

func GetBuildMode(cwd string) (BuildMode, error) {
	mode, err := ioutil.ReadFile(path.Join(cwd, defaultWorkPath, "buildMode"))
	if err != nil {
		return "", err
	}
	return BuildMode(mode), nil
}

func (a *ACBuild) rehashAndStoreOCIBlob(targetPath string, newLayer bool) error {
	layerDigestWriter := sha256.New()

	finishedWriting := false

	tmpFile, err := ioutil.TempFile(a.ContextPath, "acbuild-layer-rehashing")
	if err != nil {
		return err
	}
	defer func() {
		if !finishedWriting {
			tmpFile.Close()
		}
	}()
	combinedWriter := io.MultiWriter(layerDigestWriter, tmpFile)

	gzipWriter := gzip.NewWriter(combinedWriter)
	defer func() {
		if !finishedWriting {
			gzipWriter.Close()
		}
	}()

	diffIdWriter := sha256.New()
	tarWriter := tar.NewWriter(io.MultiWriter(diffIdWriter, gzipWriter))
	defer func() {
		if !finishedWriting {
			tarWriter.Close()
		}
	}()

	err = filepath.Walk(targetPath, util.PathWalker(tarWriter, targetPath))
	if err != nil {
		return err
	}

	tarWriter.Close()
	gzipWriter.Close()
	tmpFile.Close()

	finfo, err := os.Stat(tmpFile.Name())
	if err != nil {
		return err
	}
	fsize := finfo.Size()

	finishedWriting = true

	// See https://github.com/opencontainers/image-spec/blob/master/config.md for the difference between layer
	// digest and DiffID.
	layerDigest := hex.EncodeToString(layerDigestWriter.Sum(nil))
	diffId := hex.EncodeToString(diffIdWriter.Sum(nil))

	err = os.MkdirAll(path.Join(a.CurrentImagePath, "blobs", "sha256"), 0755)
	if err != nil {
		return err
	}

	err = os.Rename(tmpFile.Name(), path.Join(a.CurrentImagePath, "blobs", "sha256", layerDigest))
	if err != nil {
		return err
	}

	blobStorePath := path.Dir(path.Dir(targetPath))
	err = os.Rename(targetPath, path.Join(blobStorePath, "sha256", layerDigest))
	if err != nil {
		return err
	}

	var oldTopLayerHash string
	switch ociMan := a.man.(type) {
	case *oci.Image:
		if newLayer {
			// add a new top layer to the config/manifest
			err = ociMan.NewTopLayer("sha256", layerDigest, diffId, fsize)
			if err != nil {
				return err
			}
		} else {
			// update the top layer hash in the config/manifest, and remove the old
			// top layer
			oldTopLayerHash, err = ociMan.UpdateTopLayer("sha256", layerDigest, diffId, fsize)
			if err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("mismatch between build mode and manifest type?!")
	}
	if !newLayer && oldTopLayerHash != "" {
		err = os.Remove(path.Join(a.CurrentImagePath, "blobs", strings.Replace(oldTopLayerHash, ":", "/", -1)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error removing old top layer, hash %s: %v", oldTopLayerHash, err)
		}
	}

	return nil
}
