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

package appc

// SetWorkingDir sets the workingDirectory value in the untarred ACI stored at
// a.CurrentImagePath
func (m *Manifest) SetWorkingDir(dir string) error {
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	m.manifest.App.WorkingDirectory = dir
	return m.save()
}
