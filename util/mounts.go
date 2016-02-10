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

package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func IsMounted(path string) (bool, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	file, err := ioutil.ReadFile("/proc/mounts")
	if err != nil {
		return false, err
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		if len(tokens) < 2 {
			continue
		}
		if tokens[1] == absPath {
			return true, nil
		}
	}
	return false, nil
}

func MaybeUnmount(path string) error {
	_, err := os.Stat(path)
	switch {
	case os.IsNotExist(err):
		break
	case err != nil:
		return err
	default:
		mounted, err := IsMounted(path)
		if err != nil {
			return err
		}
		if mounted {
			err = syscall.Unmount(path, 0)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
