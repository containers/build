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

func (i *Image) AddAnnotation(name, value string) error {
	i.addAnnotationSaveless(name, value)
	return i.save()
}

func (i *Image) addAnnotationSaveless(name, value string) {
	if i.manifest.Annotations == nil {
		i.manifest.Annotations = make(map[string]string)
	}
	i.manifest.Annotations[name] = value
}

func (i *Image) getAnnotation(name string) string {
	if i.manifest.Annotations == nil {
		i.manifest.Annotations = make(map[string]string)
	}
	return i.manifest.Annotations[name]
}

func (i *Image) RemoveAnnotation(name string) error {
	err := i.removeAnnotationSaveless(name)
	if err != nil {
		return err
	}
	return i.save()
}

func (i *Image) removeAnnotationSaveless(name string) error {
	if i.manifest.Annotations == nil {
		i.manifest.Annotations = make(map[string]string)
	}
	_, ok := i.manifest.Annotations[name]
	if !ok {
		return fmt.Errorf("no annotation with name %q to remove", name)
	}
	delete(i.manifest.Annotations, name)
	return nil
}

func (i *Image) removeAnnotationByValSaveless(value string) error {
	fmt.Printf("removing annotation by value %q\n", value)
	if i.manifest.Annotations == nil {
		i.manifest.Annotations = make(map[string]string)
	}
	for k, v := range i.manifest.Annotations {
		if v == value {
			delete(i.manifest.Annotations, k)
			return nil
		}
	}
	return fmt.Errorf("no annotation with value %q to remove", value)
}

func (i *Image) GetAnnotations() (map[string]string, error) {
	return i.manifest.Annotations, nil
}
