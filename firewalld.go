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
	"io"

	"github.com/godbus/dbus/v5"
)

type caller interface {
	Call(ctx context.Context, c call) error
}

// Client for the Firewalld DBUS API
type Client struct {
	conn io.Closer
	dbus caller
}

const (
	objectDest = "org.fedoraproject.FirewallD1"
	objectPath = "/org/fedoraproject/FirewallD1"
)

// Opens a new connection to the system dbus and returns a connected Client for firewalld.
func Open() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		dbus: &dbusObjectWrapper{
			obj: conn.Object(objectDest, objectPath)},
	}, nil
}

const getPropertyMethod = "org.freedesktop.DBus.Properties.Get"

// Returns the Firewalld version
func (c *Client) Version(ctx context.Context) (string, error) {
	var version string
	return version, c.dbus.Call(ctx,
		newCall(getPropertyMethod, 0).
			WithArguments("org.fedoraproject.FirewallD1", "version").
			WithReturns(&version))
}

// Close disconnects from dbus
func (c *Client) Close() error {
	return c.conn.Close()
}

// dbusObjectWrapper implements the caller interface via dbus.BusObject
type dbusObjectWrapper struct {
	obj dbus.BusObject
}

var _ caller = (*dbusObjectWrapper)(nil)

func (w *dbusObjectWrapper) Call(ctx context.Context, c call) error {
	return w.obj.
		CallWithContext(ctx, c.Method, c.Flags, c.Arguments...).
		Store(c.Returns...)
}

// call is a container for all DBUS call parameters
type call struct {
	Method    string
	Flags     dbus.Flags
	Arguments []interface{}
	Returns   []interface{}
}

func newCall(method string, flags dbus.Flags) call {
	return call{
		Method: method,
		Flags:  flags,
	}
}

func (c call) WithArguments(args ...interface{}) call {
	c.Arguments = args
	return c
}

func (c call) WithReturns(returns ...interface{}) call {
	c.Returns = returns
	return c
}
