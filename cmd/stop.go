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
	"os/exec"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Cage the daemon",
	Run: func(cmd *cobra.Command, args []string) {
		disableDaemon()
		fmt.Println("kumad stopped and disabled")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}

func disableDaemon() {
	// Disable the kumad service from starting on boot
	syscmd := exec.Command("systemctl", "disable", "kumad")
	err := syscmd.Run()
	if err != nil {
		fmt.Printf("Error disabling kumad service: %v\n", err)
		return
	}

	// Stop the kumad service
	syscmd = exec.Command("systemctl", "stop", "kumad")
	err = syscmd.Run()
	if err != nil {
		fmt.Printf("Error stopping kumad service: %v\n", err)
		return
	}
}
