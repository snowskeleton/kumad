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
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Free the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		enableDaemon()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func enableDaemon() {

    // systemdConnection, _ := dbus.NewSystemdConnection()

    // listOfUnits, _ := systemdConnection.ListUnits()

    // for _, unit := range listOfUnits {
    //     fmt.Println(unit.Name)
    // }

    // systemdConnection.Close()

	// return


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
	// servicePath := "./kumad.service"
	servicePath := "/etc/systemd/system/kumad.service"
	err := os.WriteFile(servicePath, []byte(serviceContent), 0644)
	if err != nil {
		fmt.Printf("Error writing service file: %v\n", err)
		return
	}

	// Copy the kumad binary to a system-wide location
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}

	destinationPath := "/usr/local/bin/kumad"
	input, err := os.ReadFile(executablePath)
	if err != nil {
		fmt.Printf("Error reading executable: %v\n", err)
		return
	}

	err = os.WriteFile(destinationPath, input, 0755)
	if err != nil {
		fmt.Printf("Error writing executable to /usr/local/bin: %v\n", err)
		return
	}

	// Reload systemd manager configuration
	syscmd := exec.Command("systemctl", "daemon-reload")
	err = syscmd.Run()
	if err != nil {
		fmt.Printf("Error reloading systemctl daemon: %v\n", err)
		return
	}

	// Enable the kumad service to start on boot
	syscmd = exec.Command("systemctl", "enable", "kumad")
	err = syscmd.Run()
	if err != nil {
		fmt.Printf("Error enabling kumad service: %v\n", err)
		return
	}

	// Start the kumad service
	syscmd = exec.Command("systemctl", "start", "kumad")
	err = syscmd.Run()
	if err != nil {
		fmt.Printf("Error starting kumad service: %v\n", err)
		return
	}

	fmt.Println("kumad service installed and started successfully")
}
