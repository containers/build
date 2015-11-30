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

package main

import (
	"fmt"
	"strings"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/discovery"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/spf13/cobra"
)

var (
	imageId string
	labels  labellist
	size    uint
	cmdDep  = &cobra.Command{
		Use:   "dependency [command]",
		Short: "Manage dependencies",
	}
	cmdAddDep = &cobra.Command{
		Use:     "add IMAGE_NAME",
		Short:   "Add a dependency",
		Long:    "Updates the ACI to contain a dependency with the given name. If the dependency already exists, its values will be changed.",
		Example: "acbuild dependency add example.com/reduce-worker-base --label os=linux --label env=canary --size 22017258",
		Run:     runWrapper(runAddDep),
	}
	cmdRmDep = &cobra.Command{
		Use:     "remove IMAGE_NAME",
		Aliases: []string{"rm"},
		Short:   "Remove a dependency",
		Long:    "Removes the dependency with the given name from the ACI's manifest",
		Example: "acbuild dependency remove example.com/reduce-worker-base",
		Run:     runWrapper(runRmDep),
	}
)

func init() {
	cmdAcbuild.AddCommand(cmdDep)
	cmdDep.AddCommand(cmdAddDep)
	cmdDep.AddCommand(cmdRmDep)

	cmdAddDep.Flags().StringVar(&imageId, "image-id", "", "Content hash of the dependency")
	cmdAddDep.Flags().Var(&labels, "label", "Labels used for dependency matching")
	cmdAddDep.Flags().UintVar(&size, "size", 0, "The size of the image of the referenced dependency, in bytes")
}

func runAddDep(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("dependency add: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Adding dependency %q", args[0])
	}

	app, err := discovery.NewAppFromString(args[0])
	if err != nil {
		stderr("dependency add: couldn't parse dependency name: %v", err)
		return 1
	}

	appcLabels := types.Labels(labels)

	for name, value := range app.Labels {
		if _, ok := appcLabels.Get(string(name)); ok {
			stderr("multiple %s labels specified", name)
			return 1
		}
		appcLabels = append(appcLabels, types.Label{
			Name:  name,
			Value: value,
		})
	}

	var hash *types.Hash
	if imageId != "" {
		var err error
		hash, err = types.NewHash(imageId)
		if err != nil {
			stderr("dependency add: couldn't parse image ID: %v", err)
			return 1
		}
	}

	err = newACBuild().AddDependency(app.Name, hash, appcLabels, size)

	if err != nil {
		stderr("dependency add: %v", err)
		return getErrorCode(err)
	}

	return 0
}

func runRmDep(cmd *cobra.Command, args []string) (exit int) {
	if len(args) == 0 {
		cmd.Usage()
		return 1
	}
	if len(args) != 1 {
		stderr("dependency remove: incorrect number of arguments")
		return 1
	}

	if debug {
		stderr("Removing dependency %q", args[0])
	}

	err := newACBuild().RemoveDependency(args[0])

	if err != nil {
		stderr("dependency remove: %v", err)
		return getErrorCode(err)
	}

	return 0
}

type labellist []types.Label

func (ls *labellist) String() string {
	strLabels := make([]string, len(*ls))
	for i, label := range *ls {
		strLabels[i] = fmt.Sprintf("%s=%s", label.Name, label.Value)
	}
	return strings.Join(strLabels, " ")
}

func (ls *labellist) Set(input string) error {
	parts := strings.SplitN(input, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("no '=' character in %q", input)
	}
	acid, err := types.NewACIdentifier(parts[0])
	if err != nil {
		return err
	}
	*ls = append(*ls, types.Label{
		Name:  *acid,
		Value: parts[1],
	})
	return nil
}

func (ls *labellist) Type() string {
	return "Labels"
}
