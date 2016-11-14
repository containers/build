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
	"strconv"
	"strings"
)

const (
	portAnnoNamePattern  = "coreos.com/acbuild/port/%s"
	portAnnoValuePattern = "number:%d protocol:%s"
)

func (i *Image) AddPort(name, protocol string, port, count uint, socketActivated bool) error {
	if count != 1 {
		return fmt.Errorf("oci build mode does not support port counts")
	}
	if socketActivated {
		return fmt.Errorf("oci build mode does not support socket activated ports")
	}

	annoName := fmt.Sprintf(portAnnoNamePattern, name)

	if i.getAnnotation(annoName) != "" {
		return fmt.Errorf("port with name %q already exists", name)
	}

	annoValue := fmt.Sprintf(portAnnoValuePattern, port, protocol)
	i.addAnnotationSaveless(annoName, annoValue)

	if i.config.Config.ExposedPorts == nil {
		i.config.Config.ExposedPorts = make(map[string]struct{})
	}
	i.config.Config.ExposedPorts[fmt.Sprintf("%d/%s", port, protocol)] = struct{}{}

	return i.save()
}

func (i *Image) RemovePort(port string) error {
	if i.config.Config.ExposedPorts == nil {
		i.config.Config.ExposedPorts = make(map[string]struct{})
	}
	// If this port is a number or a number/protocol, check if it exists
	for key := range i.config.Config.ExposedPorts {
		tokens := strings.Split(key, "/")
		if key == port || tokens[0] == port {
			// It does exist, delete it, any related annotation, and return
			delete(i.config.Config.ExposedPorts, key)
			num, err := strconv.Atoi(tokens[0])
			if err != nil {
				return err
			}
			annoValue := fmt.Sprintf(portAnnoValuePattern, num, tokens[1])
			i.removeAnnotationByValSaveless(annoValue)
			return i.save()
		}
	}
	// If this port is a name, check for a matching annotation
	annoName := fmt.Sprintf(portAnnoNamePattern, port)
	annoValue := i.getAnnotation(annoName)
	if annoValue != "" {
		// We found an annotation, let's parse out the number and protocol
		var number uint
		var prot string
		n, err := fmt.Sscanf(annoValue, portAnnoValuePattern, &number, &prot)
		if n == 2 && err == nil {
			// The values were scanned successfully, let's see if this port exists
			str := fmt.Sprintf("%d", number) + "/" + prot
			_, ok := i.config.Config.ExposedPorts[str]
			if ok {
				// It does exist, delete it and return
				delete(i.config.Config.ExposedPorts, str)
				i.removeAnnotationSaveless(annoName)
				return i.save()
			}
		}
	}
	return fmt.Errorf("no such port %q", port)
}
