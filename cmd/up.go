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
	"os"
	"os/user"
	"time"

	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/spf13/cobra"
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Free the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		enableDaemon()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

func enableDaemon() {
	// We need to be root to edit systemd files
    currentUser, err := user.Current()
    if err != nil {
        log.Fatalf("[isRoot] Unable to get current user: %s", err)
    }
    if currentUser.Username != "root" {
		fmt.Println("Please run as root. E.g.,\n\tsudo !!")
		return
	}


	serviceContent := `[Unit]
Description=Kumad Daemon
After=network.target

[Service]
ExecStart=/usr/local/bin/kumad unattended
Restart=always
User=root

[Install]
WantedBy=multi-user.target
`

	// Write the service file to the systemd directory
	servicePath := "/etc/systemd/system/kumad.service"

	err = os.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		fmt.Printf("Error writing service file: %v\n", err)
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

	conn.ReloadContext(ctx)

	// Stop previous service
    jobID, err := conn.StopUnitContext(ctx, "kumad.service", "replace", nil)
    if err != nil {
        log.Fatalf("Failed to stop service: %v", err)
    }
    fmt.Println("Service stopped successfully. Job ID:", jobID)

    // Enable the service
    changes, _, err := conn.EnableUnitFilesContext(ctx, []string{"kumad.service"}, false, false)
	if err != nil {
        log.Fatalf("Failed to enable service: %v", err)
    }
    fmt.Println("Service enabled successfully. Changes:", changes)

    // Start the service
    jobID, err = conn.StartUnitContext(ctx, "kumad.service", "replace", nil)
    if err != nil {
        log.Fatalf("Failed to start service: %v", err)
    }
    fmt.Println("Service started successfully. Job ID:", jobID)
}
