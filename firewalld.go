/*
Copyright 2021 The routerd authors.

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

type connection interface {
	io.Closer
	Object(dest, path string) caller
}

// Opens a new connection to the system dbus and returns a connected Client for firewalld.
func Open() (*Client, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	c := &dbusConnectionWrapper{conn: conn}
	return NewClient(c), nil
}

// Client for the Firewalld D-Bus API
type Client struct {
	conn   connection
	main   caller
	config *ConfigClient
}

const (
	// Firewalld D-Bus destination
	dbusDest = "org.fedoraproject.FirewallD1"

	mainPath   = "/org/fedoraproject/FirewallD1"
	configPath = "/org/fedoraproject/FirewallD1/config"
)

func NewClient(conn connection) *Client {
	return &Client{
		conn: conn,
		main: conn.Object(dbusDest, mainPath),

		config: NewConfigClient(conn),
	}
}

// Config returns a client for working on firewalld persistant configuration.
func (c *Client) Config() *ConfigClient {
	return c.config
}

const getPropertyMethod = "org.freedesktop.DBus.Properties.Get"

// Returns the Firewalld version
func (c *Client) Version(ctx context.Context) (string, error) {
	var version string
	return version, c.main.Call(ctx,
		newCall(getPropertyMethod, 0).
			WithArguments("org.fedoraproject.FirewallD1", "version").
			WithReturns(&version))
}

const reloadMethod = "org.fedoraproject.FirewallD1.reload"

func (c *Client) Reload(ctx context.Context) error {
	return c.main.Call(ctx,
		newCall(reloadMethod, 0))
}

// Close disconnects from dbus
func (c *Client) Close() error {
	return c.conn.Close()
}

type dbusConnectionWrapper struct {
	conn *dbus.Conn
}

var _ connection = (*dbusConnectionWrapper)(nil)

func (w *dbusConnectionWrapper) Close() error {
	return w.conn.Close()
}

func (w *dbusConnectionWrapper) Object(dest, path string) caller {
	return &dbusObjectWrapper{
		obj: w.conn.Object(dest, dbus.ObjectPath(path))}
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
