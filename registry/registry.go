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

// The registry package exists to manage ACIs for acbuild. The main difference
// between this package and rkt's store is that this package is optimised for
// many separate calls into the current ACI, as opposed to having a tinier
// footprint. When an ACI is fetched, it is immediately rendered onto the
// filesystem, so that when acbuild's run command is invoked many times there's
// no waiting for files to be untarred or uncompressed.
package registry

import (
	"crypto/sha512"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"

	"github.com/appc/acbuild/util"
)

var (
	hashPrefix = "sha512-"
	lenHash    = sha512.Size       // raw byte size
	lenHashKey = (lenHash / 2) * 2 // half length, in hex characters
	lenKey     = len(hashPrefix) + lenHashKey
	minlenKey  = len(hashPrefix) + 2 // at least sha512-aa

	ErrNotFound = fmt.Errorf("ACI not in registry")
)

type Registry struct {
	DepStoreTarPath      string
	DepStoreExpandedPath string
	Insecure             bool
	Debug                bool
}

// Read the ACI contents stream given the key. Use ResolveKey to
// convert an image ID to the relative provider's key.
func (r Registry) ReadStream(key string) (io.ReadCloser, error) {
	return os.Open(path.Join(r.DepStoreTarPath, key))
}

// Converts an image ID to the, if existent, key under which the
// ACI is known to the provider
func (r Registry) ResolveKey(key string) (string, error) {
	if !strings.HasPrefix(key, hashPrefix) {
		return "", fmt.Errorf("wrong key prefix")
	}
	if len(key) < minlenKey {
		return "", fmt.Errorf("key too short")
	}
	if len(key) > lenKey {
		key = key[:lenKey]
	}

	files, err := ioutil.ReadDir(r.DepStoreTarPath)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if len(file.Name()) >= len(key) && key == file.Name()[:len(key)] {
			return key, nil
		}
	}
	return "", fmt.Errorf("key not found in registry")
}

// Converts a Hash to the provider's key
func (r Registry) HashToKey(h hash.Hash) string {
	s := h.Sum(nil)
	if len(s) != sha512.Size {
		fmt.Fprintln(os.Stderr, "bad hash passed to the registry")
		// Return a nonsensical key that won't resolve to anything
		return "libacb-bad-registry-key"
	}
	return fmt.Sprintf("%s%x", hashPrefix, s)
}

// Returns the manifest for the ACI with the given key
func (r Registry) GetImageManifest(key string) (*schema.ImageManifest, error) {
	return util.GetManifest(path.Join(r.DepStoreExpandedPath, key))
}

// Returns the key for the ACI with the given name and labels
func (r Registry) GetACI(name types.ACIdentifier, labels types.Labels) (string, error) {
	files, err := ioutil.ReadDir(r.DepStoreExpandedPath)
	if err != nil {
		return "", err
	}
nextkey:
	for _, file := range files {
		man, err := util.GetManifest(path.Join(r.DepStoreExpandedPath, file.Name()))
		if err != nil {
			return "", err
		}
		if man.Name != name {
			continue
		}
		for _, l := range labels {
			val, ok := man.Labels.Get(l.Name.String())
			if !ok {
				continue
			}
			if val != l.Value {
				continue nextkey
			}
		}
		return file.Name(), nil
	}
	return "", ErrNotFound
}
