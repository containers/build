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
	"fmt"
)

const (
	mountAnnoNamePattern  = "coreos.com/acbuild/mount/%s"
	mountAnnoValuePattern = "path:%s"
)

func (i *Image) AddMount(name, path string, readOnly bool) error {
	if readOnly {
		return fmt.Errorf("oci build mode does not support read only mounts")
	}
	if i.config.Config.Volumes == nil {
		i.config.Config.Volumes = make(map[string]struct{})
	}
	i.config.Config.Volumes[path] = struct{}{}

	annoName := fmt.Sprintf(mountAnnoNamePattern, name)

	if i.getAnnotation(annoName) != "" {
		return fmt.Errorf("mount with name %q already exists", name)
	}

	annoValue := fmt.Sprintf(mountAnnoValuePattern, path)
	i.addAnnotationSaveless(annoName, annoValue)

	return i.save()
}

func (i *Image) RemoveMount(mount string) error {
	if i.config.Config.Volumes == nil {
		i.config.Config.Volumes = make(map[string]struct{})
	}
	// If this mount is a path, check if it exists
	_, ok := i.config.Config.Volumes[mount]
	if ok {
		// It does! Great! Delete it, any related annotation, and return.
		delete(i.config.Config.Volumes, mount)
		annoValue := fmt.Sprintf(mountAnnoValuePattern, mount)
		i.removeAnnotationByValSaveless(annoValue)
		return i.save()
	}
	// If this mount is a name, check for a matching annotation
	annoName := fmt.Sprintf(mountAnnoNamePattern, mount)
	annoValue := i.getAnnotation(annoName)
	if annoValue != "" {
		// We found an annotation! Let's scan out the path
		var path string
		n, err := fmt.Sscanf(annoValue, mountAnnoValuePattern, &path)
		if n == 1 && err == nil {
			// The path was scanned successfully, let's see if it exists
			_, ok := i.config.Config.Volumes[path]
			if ok {
				// It does! Great! Delete it, any related annotation, and return.
				delete(i.config.Config.Volumes, path)
				i.removeAnnotationSaveless(annoName)
				return i.save()
			}
		}
	}
	// Couldn't find a matching mount :(
	return fmt.Errorf("no such mount: %s", mount)
}
