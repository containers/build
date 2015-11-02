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

package tests

import (
	"strconv"
	"testing"

	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema"
	"github.com/appc/acbuild/Godeps/_workspace/src/github.com/appc/spec/schema/types"
)

const (
	portName          = "http"
	portProtocol      = "tcp"
	portNumber   uint = 8080
	portCount    uint = 2

	portName2     = "snmp"
	portProtocol2 = "udp"
	portNumber2   = 161
)

func manWithPorts(ports []types.Port) schema.ImageManifest {
	return schema.ImageManifest{
		ACKind:    schema.ImageManifestKind,
		ACVersion: schema.AppContainerVersion,
		Name:      *types.MustACIdentifier("acbuild-unnamed"),
		App: &types.App{
			Exec:  nil,
			User:  "0",
			Group: "0",
			Ports: ports,
		},
		Labels: systemLabels,
	}
}

func TestAddPort(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	ports := []types.Port{
		types.Port{
			Name:     *types.MustACName(portName),
			Protocol: portProtocol,
			Port:     portNumber,
			Count:    1,
		},
	}

	checkManifest(t, workingDir, manWithPorts(ports))
	checkEmptyRootfs(t, workingDir)
}

func TestAddPortWithCount(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)),
		"--count", strconv.Itoa(int(portCount)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	ports := []types.Port{
		types.Port{
			Name:     *types.MustACName(portName),
			Protocol: portProtocol,
			Port:     portNumber,
			Count:    portCount,
		},
	}

	checkManifest(t, workingDir, manWithPorts(ports))
	checkEmptyRootfs(t, workingDir)
}

func TestAddPortWithSocketActivated(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)),
		"--socket-activated")
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	ports := []types.Port{
		types.Port{
			Name:            *types.MustACName(portName),
			Protocol:        portProtocol,
			Port:            portNumber,
			Count:           1,
			SocketActivated: true,
		},
	}

	checkManifest(t, workingDir, manWithPorts(ports))
	checkEmptyRootfs(t, workingDir)
}

func TestAddNegativePort(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	// The "--" is required to prevent cobra from parsing the "-1" as a flag
	err := runACBuild(workingDir, "port", "add", portName, portProtocol, "--", "-1")
	if err == nil {
		t.Fatalf("port add didn't return an error when asked to add a port with a negative number")
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}

func TestAddPortThatsTooHigh(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, "65536")
	if err == nil {
		t.Fatalf("port add didn't return an error when asked to add a port with a number > 65535")
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}

func TestAddTwoPorts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "port", "add", portName2, portProtocol2, strconv.Itoa(int(portNumber2)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	ports := []types.Port{
		types.Port{
			Name:     *types.MustACName(portName),
			Protocol: portProtocol,
			Port:     portNumber,
			Count:    1,
		},
		types.Port{
			Name:     *types.MustACName(portName2),
			Protocol: portProtocol2,
			Port:     portNumber2,
			Count:    1,
		},
	}

	checkManifest(t, workingDir, manWithPorts(ports))
	checkEmptyRootfs(t, workingDir)
}

func TestAddRmPorts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "port", "rm", portName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	checkManifest(t, workingDir, emptyManifestWithApp)
	checkEmptyRootfs(t, workingDir)
}

func TestAddAddRmPorts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "add", portName, portProtocol, strconv.Itoa(int(portNumber)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "port", "add", portName2, portProtocol2, strconv.Itoa(int(portNumber2)))
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	err = runACBuild(workingDir, "port", "rm", portName)
	if err != nil {
		t.Fatalf("%v\n", err)
	}

	ports := []types.Port{
		types.Port{
			Name:     *types.MustACName(portName2),
			Protocol: portProtocol2,
			Port:     portNumber2,
			Count:    1,
		},
	}

	checkManifest(t, workingDir, manWithPorts(ports))
	checkEmptyRootfs(t, workingDir)
}

func TestRmNonexistentPorts(t *testing.T) {
	workingDir := setUpTest(t)
	defer cleanUpTest(workingDir)

	err := runACBuild(workingDir, "port", "remove", portName)
	switch {
	case err == nil:
		t.Fatalf("port remove didn't return an error when asked to remove nonexistent port")
	case err.exitCode == 2:
		return
	default:
		t.Fatalf("error occurred when running port remove:\n%v", err)
	}

	checkManifest(t, workingDir, emptyManifest)
	checkEmptyRootfs(t, workingDir)
}
