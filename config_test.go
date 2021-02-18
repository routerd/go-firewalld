/*
Copyright 2021 The routerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package firewalld

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func configClientSetup() (
	configPathCaller *callerMock,
	conn *connectionMock,
	c *ConfigClient,
) {
	configPathCaller = &callerMock{}

	conn = &connectionMock{}
	conn.On("Object", dbusDest, configPath).Return(configPathCaller)

	c = NewConfigClient(conn)
	return
}

func TestConfigClient_GetZoneNames(t *testing.T) {
	response := []string{"FedoraServer", "dmz", "drop"}

	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*[]string)
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	zones, err := c.GetZoneNames(ctx)
	require.NoError(t, err)

	assert.Equal(t, response, zones)
}

func TestConfigClient_GetServiceNames(t *testing.T) {
	response := []string{"ssh", "samba-client", "dhcpv6-client"}

	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*[]string)
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	zones, err := c.GetServiceNames(ctx)
	require.NoError(t, err)

	assert.Equal(t, response, zones)
}

func TestConfigClient_ListZones(t *testing.T) {
	response := []string{
		"/org/fedoraproject/FirewallD1/config/zone/0",
		"/org/fedoraproject/FirewallD1/config/zone/1",
		"/org/fedoraproject/FirewallD1/config/zone/2",
	}

	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*[]string)
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	zones, err := c.ListZones(ctx)
	require.NoError(t, err)

	assert.Equal(t, response, zones)
}

func TestConfigClient_GetZoneByName(t *testing.T) {
	response := "/org/fedoraproject/FirewallD1/config/zone/0"

	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*string)
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	zonePath, err := c.GetZoneByName(ctx, "test")
	require.NoError(t, err)

	assert.Equal(t, response, zonePath)
}

func TestConfigClient_GetServiceByName(t *testing.T) {
	response := "/org/fedoraproject/FirewallD1/config/service/138"

	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*string)
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	servicePath, err := c.GetServiceByName(ctx, "ssh")
	require.NoError(t, err)

	assert.Equal(t, response, servicePath)
}

func TestConfigClient_RemoveZone(t *testing.T) {
	const path = "/org/fedoraproject/FirewallD1/config/zone/0"

	configPathCaller, conn, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.MatchedBy(func(o interface{}) bool {
			c := o.(call)
			return c.Method == configGetZoneByNameMethod
		})).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*string)
			*s = path
		}).
		Return(nil)

	zoneObjectCaller := &callerMock{}
	conn.
		On("Object", dbusDest, path).
		Return(zoneObjectCaller)
	zoneObjectCaller.
		On("Call", mock.Anything, mock.Anything).
		Return(nil)

	ctx := context.Background()

	err := c.RemoveZone(ctx, "ssh")
	require.NoError(t, err)
}

func TestConfigClient_AddZone(t *testing.T) {
	configPathCaller, _, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.Anything).
		Return(nil)

	ctx := context.Background()

	err := c.AddZone(ctx, "ssh", ZoneSettings{})
	require.NoError(t, err)
}

func TestClient_GetZoneSettings(t *testing.T) {
	expected := ZoneSettings{
		Version:     "",
		Name:        "Fedora Workstation",
		Description: "Unsolicited incoming network packets are rejected from port 1 to 1024, except for select network services. Incoming packets that are related to outgoing network connections are accepted. Outgoing network connections are allowed.",
		Target:      "default",
		Services:    []string{"dhcpv6-client", "ssh", "samba-client", "mdns"},
		Ports: []Port{
			{Port: "1025-65535", Protocol: "udp"},
			{Port: "1025-65535", Protocol: "tcp"},
		},
		// ICMPBlocks:      []string{},
		Masquerade: false,
		ForwardPorts: []ForwardPort{
			{Port: "22", Protocol: "tcp", ToAddress: "192.0.2.55", ToPort: "22"},
		},
		Interfaces: []string{"wlp4s0", "tun0", "ens1u2u1u2"},
		// SourceAddresses: []string{},
		// RichRules:       []string{},
		// Protocols:       []string{},
		SourcePorts: []Port{
			{Port: "1025-65535", Protocol: "udp"},
			{Port: "1025-65535", Protocol: "tcp"},
		},
	}

	var response = []interface{}{
		expected.Version,
		expected.Name,
		expected.Description,
		false, // Unused
		expected.Target,
		expected.Services,
		[][]interface{}{ // ports
			{"1025-65535", "udp"},
			{"1025-65535", "tcp"},
		},
		expected.ICMPBlocks,
		expected.Masquerade,
		[][]interface{}{ // forward ports
			{"22", "tcp", "22", "192.0.2.55"},
		},
		expected.Interfaces,
		expected.SourceAddresses,
		expected.RichRules,
		expected.Protocols,
		[][]interface{}{ // source ports
			{"1025-65535", "udp"},
			{"1025-65535", "tcp"},
		},
		false,
	}

	const path = "/org/fedoraproject/FirewallD1/config/zone/0"

	configPathCaller, conn, c := configClientSetup()
	configPathCaller.
		On("Call", mock.Anything, mock.MatchedBy(func(o interface{}) bool {
			c := o.(call)
			return c.Method == configGetZoneByNameMethod
		})).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*string)
			*s = path
		}).
		Return(nil)

	zoneObjectCaller := &callerMock{}
	conn.
		On("Object", dbusDest, path).
		Return(zoneObjectCaller)
	zoneObjectCaller.
		On("Call", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			c := args.Get(1).(call)
			s := c.Returns[0].(*[]interface{})
			*s = response
		}).
		Return(nil)

	ctx := context.Background()

	settings, err := c.GetZoneSettings(ctx, "ssh")
	require.NoError(t, err)

	assert.Equal(t, expected, settings)
}
