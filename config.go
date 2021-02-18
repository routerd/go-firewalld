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
)

// Client for Firewalld org.fedoraproject.FirewallD1.config.
// Methods manipulate the persistent firewalld configuration.
type ConfigClient struct {
	conn       connection
	configPath caller
}

func NewConfigClient(conn connection) *ConfigClient {
	return &ConfigClient{
		conn:       conn,
		configPath: conn.Object(dbusDest, configPath),
	}
}

const configGetZoneNamesMethod = "org.fedoraproject.FirewallD1.config.getZoneNames"

// Return list of zone names (permanent configuration).
func (c *ConfigClient) GetZoneNames(
	ctx context.Context) ([]string, error) {
	var zoneNames []string
	return zoneNames, c.configPath.Call(ctx,
		newCall(configGetZoneNamesMethod, 0).
			WithReturns(&zoneNames))
}

const configGetServiceNamesMethod = "org.fedoraproject.FirewallD1.config.getServiceNames"

// Return list of service names (permanent configuration).
func (c *ConfigClient) GetServiceNames(
	ctx context.Context) ([]string, error) {
	var serviceNames []string
	return serviceNames, c.configPath.Call(ctx,
		newCall(configGetServiceNamesMethod, 0).
			WithReturns(&serviceNames))
}

const configListZonesMethod = "org.fedoraproject.FirewallD1.config.listZones"

// List object paths of zones known to permanent environment.
func (c *ConfigClient) ListZones(
	ctx context.Context) (zonePaths []string, err error) {
	return zonePaths, c.configPath.Call(ctx,
		newCall(configListZonesMethod, 0).
			WithReturns(&zonePaths))
}

const configGetZoneByNameMethod = "org.fedoraproject.FirewallD1.config.getZoneByName"

// Return object path (permanent configuration) of zone with given name.
func (c *ConfigClient) GetZoneByName(
	ctx context.Context, zoneName string) (zonePath string, err error) {
	return zonePath, c.configPath.Call(ctx,
		newCall(configGetZoneByNameMethod, 0).
			WithArguments(zoneName).
			WithReturns(&zonePath))
}

const configGetServiceByNameMethod = "org.fedoraproject.FirewallD1.config.getServiceByName"

// Return object path (permanent configuration) of service with given name.
func (c *ConfigClient) GetServiceByName(
	ctx context.Context, serviceName string) (servicePath string, err error) {
	return servicePath, c.configPath.Call(ctx,
		newCall(configGetServiceByNameMethod, 0).
			WithArguments(serviceName).
			WithReturns(&servicePath))
}

// Zone instance object interface
const configZoneRemoveMethod = "org.fedoraproject.FirewallD1.config.zone.remove"

// Remove zone with given settings into permanent configuration.
func (c *ConfigClient) RemoveZone(
	ctx context.Context, zoneName string) error {
	path, err := c.GetZoneByName(ctx, zoneName)
	if err != nil {
		return err
	}
	return c.conn.Object(dbusDest, path).
		Call(ctx, newCall(configZoneRemoveMethod, 0))
}

const addZoneMethod = "org.fedoraproject.FirewallD1.config.addZone"

// Add zone with given settings into permanent configuration.
func (c *ConfigClient) AddZone(
	ctx context.Context, zoneName string, settings ZoneSettings) error {
	var z interface{}
	return c.configPath.Call(ctx,
		newCall(addZoneMethod, 0).
			WithArguments(zoneName, settings.ToSlice()).
			WithReturns(&z))
}

const configZoneGetSettingsMethod = "org.fedoraproject.FirewallD1.config.zone.getSettings"

// Return permanent settings of given zone.
func (c *ConfigClient) GetZoneSettings(
	ctx context.Context, zoneName string) (ZoneSettings, error) {
	path, err := c.GetZoneByName(ctx, zoneName)
	if err != nil {
		return ZoneSettings{}, err
	}

	var zoneSettings []interface{}
	err = c.conn.Object(dbusDest, path).
		Call(ctx,
			newCall(configZoneGetSettingsMethod, 0).
				WithReturns(&zoneSettings))
	if err != nil {
		return ZoneSettings{}, err
	}

	return ZoneSettingsFromSlice(zoneSettings), nil
}
