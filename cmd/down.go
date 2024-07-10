/*
Copyright © 2024 Isaac Lyons <isaac@snowskeleton.net>

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
package cmd

import (
	"context"
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/spf13/cobra"
)

// downCmd represents the stop command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Cage the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		disableDaemon()
		fmt.Println("kumad stopped and disabled")
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}

func disableDaemon() {
	// We need to be root to edit systemd files
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf("[isRoot] Unable to get current user: %s", err)
    }
    if currentUser.Username != "root" {
		fmt.Println("Please run as root. E.g.,\n\tsudo !!")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Connect to systemd
    conn, err := dbus.NewSystemdConnectionContext(ctx)
    if err != nil {
        log.Fatalf("Failed to connect to system bus: %v", err)
    }
    defer conn.Close()

    // Disable the service
    changes, err := conn.DisableUnitFilesContext(ctx, []string{"kumad.service"}, false)
	if err != nil {
        log.Fatalf("Failed to disable service: %v", err)
    }
    fmt.Println("Service disabled successfully. Changes:", changes)

    // Stop the service
    jobID, err := conn.StopUnitContext(ctx, "kumad.service", "replace", nil)
    if err != nil {
        log.Fatalf("Failed to stop service: %v", err)
    }
    fmt.Println("Service stopped successfully. Job ID:", jobID)
}
