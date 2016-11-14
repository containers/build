// Copyright 2016 The appc Authors
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

package oci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/containers/build/util"

	//specs "github.com/opencontainers/image-spec/specs-go"
	ociImage "github.com/opencontainers/image-spec/specs-go/v1"
)

var (
	// ErrNotFound is returned when acbuild is asked to remove an element from a
	// list and the element is not present in the list
	ErrNotFound = fmt.Errorf("element to be removed does not exist in this image")
)

// Manifest is a struct with an open handle to a manifest that it can manipulate
type Image struct {
	ociPath  string
	refName  string
	config   ociImage.Image
	manifest ociImage.Manifest
	ref      ociImage.Descriptor
}

func LoadImage(ociPath string) (*Image, error) {
	i := &Image{
		ociPath: ociPath,
		refName: "latest",
	}

	refDir := path.Join(ociPath, "refs")
	blobDir := path.Join(ociPath, "blobs")

	// Look for refs
	refFileInfos, err := ioutil.ReadDir(refDir)
	if err != nil {
		return nil, err
	}
	if len(refFileInfos) == 0 {
		return nil, fmt.Errorf("no refs found in image")
	}
	// We need to pick a ref, if there's more than one we don't know which one
	// the user wishes to modify. Let's just pick the first one.
	i.refName = path.Base(refFileInfos[0].Name())

	// Open the ref file, read it, unmarshal it, and parse the manifest's
	// hash
	refFile, err := os.OpenFile(path.Join(refDir, i.refName), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer refFile.Close()
	refBlob, err := ioutil.ReadAll(refFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(refBlob, &i.ref)
	if err != nil {
		return nil, err
	}
	manifestHashAlgo, manifestHash, err := splitHash(i.ref.Digest)
	if err != nil {
		return nil, err
	}

	// Open the manifest, read it, unmarshal it, and parse the config's hash
	manifestFile, err := os.OpenFile(path.Join(blobDir, manifestHashAlgo, manifestHash), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer manifestFile.Close()
	manifestBlob, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(manifestBlob, &i.manifest)
	if err != nil {
		return nil, err
	}
	configHashAlgo, configHash, err := splitHash(i.manifest.Config.Digest)
	if err != nil {
		return nil, err
	}

	// Open the config, read it, unmarshal it
	configFile, err := os.OpenFile(path.Join(blobDir, configHashAlgo, configHash), os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	configBlob, err := ioutil.ReadAll(configFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configBlob, &i.config)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func splitHash(hash string) (string, string, error) {
	pieces := strings.Split(hash, ":")
	if len(pieces) != 2 {
		return "", "", fmt.Errorf("error parsing digest in image ref: %q", hash)
	}
	return pieces[0], pieces[1], nil
}

func (i *Image) Close() error {
	return nil
}

func (i *Image) save() error {
	// Remove the old config
	oldConfigHashAlgo, oldConfigHash, err := splitHash(i.manifest.Config.Digest)
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(i.ociPath, "blobs", oldConfigHashAlgo, oldConfigHash))
	if err != nil {
		return err
	}
	// Save the new config
	configHashAlgo, configHash, configSize, err := util.MarshalHashAndWrite(i.ociPath, i.config)
	if err != nil {
		return err
	}
	i.manifest.Config.Digest = configHashAlgo + ":" + configHash
	i.manifest.Config.Size = int64(configSize)

	// Remove the old manifest
	oldManifestHashAlgo, oldManifestHash, err := splitHash(i.ref.Digest)
	if err != nil {
		return err
	}
	err = os.Remove(path.Join(i.ociPath, "blobs", oldManifestHashAlgo, oldManifestHash))
	if err != nil {
		return err
	}
	// Save the new manifest
	manifestHashAlgo, manifestHash, manifestSize, err := util.MarshalHashAndWrite(i.ociPath, i.manifest)
	if err != nil {
		return err
	}
	i.ref.Digest = manifestHashAlgo + ":" + manifestHash
	i.ref.Size = int64(manifestSize)

	// Remove any old refs
	err = os.RemoveAll(path.Join(i.ociPath, "refs"))
	if err != nil {
		return err
	}
	err = os.Mkdir(path.Join(i.ociPath, "refs"), 0755)
	if err != nil {
		return err
	}

	// Write the ref
	refBlob, err := json.Marshal(i.ref)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(i.ociPath, "refs", i.refName), refBlob, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (i *Image) GetConfig() ociImage.Image {
	return i.config
}

func (i *Image) GetManifest() ociImage.Manifest {
	return i.manifest
}

func (i *Image) GetRef() ociImage.Descriptor {
	return i.ref
}

func (i *Image) GetLayerHashes() []string {
	return i.config.RootFS.DiffIDs
}

func (i *Image) Print(prettyPrint, printConfig bool) error {
	var configblob []byte
	var err error
	var toPrint interface{}
	if printConfig {
		toPrint = i.config
	} else {
		toPrint = i.manifest
	}
	if prettyPrint {
		configblob, err = json.MarshalIndent(toPrint, "", "    ")
	} else {
		configblob, err = json.Marshal(toPrint)
	}
	if err != nil {
		return err
	}
	fmt.Println(string(configblob))
	return nil
}

func (i *Image) UpdateTopLayerHash(hashAlgo, newHash string, size int64) (string, error) {
	var oldTopLayerHash string
	hashStr := hashAlgo + ":" + newHash
	if len(i.config.RootFS.DiffIDs) == 0 {
		i.config.RootFS = ociImage.RootFS{
			Type:    "layers",
			DiffIDs: []string{hashStr},
		}
	} else {
		oldTopLayerHash = i.config.RootFS.DiffIDs[len(i.config.RootFS.DiffIDs)-1]
		i.config.RootFS.DiffIDs[len(i.config.RootFS.DiffIDs)-1] = hashStr
	}

	layerDescriptor :=
		ociImage.Descriptor{
			MediaType: "MediaTypeImageLayer",
			Digest:    hashStr,
			Size:      size,
		}
	if len(i.manifest.Layers) == 0 {
		i.manifest.Layers = []ociImage.Descriptor{layerDescriptor}
	} else {
		i.manifest.Layers[len(i.manifest.Layers)-1] = layerDescriptor
	}

	return oldTopLayerHash, i.save()
}

func (i *Image) NewTopLayer(hashAlgo, newHash string, size int64) error {
	hashStr := hashAlgo + ":" + newHash
	if len(i.config.RootFS.DiffIDs) == 0 {
		i.config.RootFS = ociImage.RootFS{
			Type:    "layers",
			DiffIDs: []string{hashStr},
		}
	} else {
		i.config.RootFS.DiffIDs = append(i.config.RootFS.DiffIDs, hashStr)
	}

	layerDescriptor :=
		ociImage.Descriptor{
			MediaType: "MediaTypeImageLayer",
			Digest:    hashStr,
			Size:      size,
		}
	if len(i.manifest.Layers) == 0 {
		i.manifest.Layers = []ociImage.Descriptor{layerDescriptor}
	} else {
		i.manifest.Layers = append(i.manifest.Layers, layerDescriptor)
	}

	return i.save()
}
