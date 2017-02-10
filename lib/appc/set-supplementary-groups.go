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

// SetSuppGroups sets the groups the pod will run as in the untarred ACI
// stored at a.CurrentACIPath.
func (m *Manifest) SetSuppGroups(groups []int) error {
	if m.manifest.App == nil {
		m.manifest.App = newManifestApp()
	}
	m.manifest.App.SupplementaryGIDs = groups
	return m.save()
}
