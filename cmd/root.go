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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var push_url string
var push_interval int
var verbose bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kumad",
	Short: "A simple agent that pings an Uptime Kuma Push monitor",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Output to console")

	rootCmd.PersistentFlags().StringVarP(&push_url, "push_url", "u", "", "The URL of the endpoint to ping")
	rootCmd.MarkFlagRequired("push_url")
	viper.BindPFlag("push_url", rootCmd.PersistentFlags().Lookup("push_url"))

	rootCmd.PersistentFlags().IntVarP(&push_interval, "push_interval", "i", 50, "The interval at which to ping")
	viper.BindPFlag("push_interval", rootCmd.PersistentFlags().Lookup("push_interval"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc/")
		viper.SetConfigType("yaml")
		viper.SetConfigName("kumad")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if rootCmd.PersistentFlags().Lookup("push_interval").Changed || viper.Get("push_interval") == nil {
		viper.Set("push_interval", push_interval)
	}
	if rootCmd.PersistentFlags().Lookup("push_url").Changed || viper.Get("push_url") == nil {
		viper.Set("push_url", push_url)
	}
	if viper.Get("push_url") == "" || viper.Get("push_url") == nil {
		fmt.Println("Please provide a value for push_url")
		os.Exit(1)
	}
	fmt.Println("Pinging URL:", viper.Get("push_url"))
	fmt.Println("Ping interval:", viper.Get("push_interval"), "seconds")
	viper.WriteConfig()
}
