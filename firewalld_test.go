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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	t.Run("Close", func(t *testing.T) {
		connection := &connectionMock{}
		connection.On("Close").Return(nil)

		c := &Client{
			conn: connection,
		}

		require.NoError(t, c.Close())
		connection.AssertCalled(t, "Close")
	})

	t.Run("Version", func(t *testing.T) {
		const response = "0.8.6"

		caller := &callerMock{}
		caller.
			On("Call",
				mock.Anything,
				mock.MatchedBy(func(c call) bool {
					return c.Method == getPropertyMethod
				})).
			Run(func(args mock.Arguments) {
				c := args.Get(1).(call)
				s := c.Returns[0].(*string)
				*s = response
			}).
			Return(nil)

		c := &Client{
			main: caller,
		}

		ctx := context.Background()
		zone, err := c.Version(ctx)
		require.NoError(t, err)

		assert.Equal(t, response, zone)
	})
}

type callerMock struct {
	mock.Mock
}

var _ caller = (*callerMock)(nil)

func (m *callerMock) Call(ctx context.Context, c call) error {
	args := m.Called(ctx, c)
	err, _ := args.Error(0).(error)
	return err
}

type connectionMock struct {
	mock.Mock
}

var _ io.Closer = (*connectionMock)(nil)

func (m *connectionMock) Close() error {
	args := m.Called()
	err, _ := args.Error(0).(error)
	return err
}

func (m *connectionMock) Object(dest, path string) caller {
	args := m.Called(dest, path)
	c, _ := args.Get(0).(caller)
	return c
}
