// Copyright 2016 The acbuild Authors
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

// SetExec sets the exec command for the untarred ACI stored at
// a.CurrentImagePath.
func (i *Image) SetExec(cmd []string) error {
	if len(cmd) > 0 {
		i.config.Config.Entrypoint = cmd[:1]
	}
	if len(cmd) > 1 {
		i.config.Config.Cmd = cmd[1:]
	}
	return i.save()
}
