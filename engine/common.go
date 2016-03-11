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

package engine

import (
	"github.com/appc/spec/schema/types"
)

// Engine is an interface which is accepted by lib.Run, and used to perform the
// actual execution of a binary inside the container.
type Engine interface {
	// Run executes a command inside a container. command is the path to the
	// binary (post-chroot) to exec, args is the arguments to pass to the
	// binary, environment is the set of environment variables to set for the
	// binary, chroot is the path on the host where the container's root
	// filesystem exists, and workingDir specifies the path inside the
	// container that should be the current working directory for the binary.
	// If workingDir is "", the default should be "/".
	Run(command string, args []string, environment types.Environment, chroot, workingDir string) error
}
