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

package main

import (
	"context"
	"fmt"
	"os"

	"routerd.net/go-firewalld"
)

func main() {
	client, err := firewalld.Open()
	exitOnErr(err)
	defer client.Close()

	ctx := context.Background()

	zones, err := client.Config().GetZoneNames(ctx)
	exitOnErr(err)
	fmt.Println("Initial Zones:", zones)

	err = client.Config().AddZone(ctx, "test", firewalld.ZoneSettings{
		Target: "default",
	})
	exitOnErr(err)
	fmt.Println("Added zone 'test'")

	zones, err = client.Config().GetZoneNames(ctx)
	exitOnErr(err)
	fmt.Println("Updated Zones:", zones)

	err = client.Config().RemoveZone(ctx, "test")
	exitOnErr(err)
	fmt.Println("Removed zone 'test'")
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "boom: ", err)
		os.Exit(1)
	}
}
