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
	"strings"
)

func (i *Image) AddEnv(name, value string) error {
	i.removeFromEnv(name)
	i.config.Config.Env = append(i.config.Config.Env, name+"="+value)
	return i.save()

}

func (i *Image) RemoveEnv(name string) error {
	err := i.removeFromEnv(name)
	if err != nil {
		return err
	}
	return i.save()
}

func (i *Image) removeFromEnv(name string) error {
	tmpEnv := make([]string, len(i.config.Config.Env))
	copy(tmpEnv, i.config.Config.Env)
	foundOne := false
	for j := len(tmpEnv) - 1; j >= 0; j-- {
		varParts := strings.SplitN(tmpEnv[j], "=", 2)
		if len(varParts) != 2 {
			return fmt.Errorf("invalid environment variable in config: %q", tmpEnv[j])
		}
		varName := varParts[0]
		if varName == name {
			foundOne = true
			tmpEnv = append(
				tmpEnv[:j],
				tmpEnv[j+1:]...,
			)
		}
	}
	if !foundOne {
		return ErrNotFound
	}
	i.config.Config.Env = tmpEnv
	return nil
}
