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

import (
	"fmt"
	"strings"
)

// SetUser will set the user (username or UID) the app in this container will
// run as
func (i *Image) SetUser(user string) error {
	if user == "" {
		return fmt.Errorf("user cannot be empty")
	}
	if strings.Contains(user, ":") {
		return fmt.Errorf("user cannot contain a ':' character")
	}
	if !strings.Contains(i.config.Config.User, ":") {
		i.config.Config.User = user
	} else {
		tokens := strings.Split(i.config.Config.User, ":")
		if len(tokens) != 2 {
			return fmt.Errorf("something has gone horribly wrong setting the user")
		}
		i.config.Config.User = user + ":" + tokens[1]
	}
	return i.save()
}
