// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	_ "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
	"strings"
	"time"
)

var cfgFile string
var registryHost string
var registryPort string

var baseRegistryURL string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "topview",
	Short: "Topology Viewer for Mashery",
	Long:  `A command-line application to view the Mashery topology`,
	//Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.topview.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	//if cfgFile != "" {
	//	// Use config file from the flag.
	//	viper.SetConfigFile(cfgFile)
	//} else {
	//	// Find home directory.
	//	home, err := homedir.Dir()
	//	if err != nil {
	//		fmt.Println(err)
	//		os.Exit(1)
	//	}
	//
	//	// Search config in home directory with name ".topview" (without extension).
	//	viper.AddConfigPath(home)
	//	viper.SetConfigName(".topview")
	//}
	//
	//viper.AutomaticEnv() // read in environment variables that match
	//
	//// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	fmt.Println("Using config file:", viper.ConfigFileUsed())
	//}
	viper.SetEnvPrefix("TMGC")
	viper.BindEnv("REGISTRY_HOST")
	viper.BindEnv("REGISTRY_PORT")

	viper.SetDefault("REGISTRY_HOST", "localhost")
	viper.SetDefault("REGISTRY_PORT", "1080")

	registryHost = viper.Get("REGISTRY_HOST").(string)
	registryPort = viper.Get("REGISTRY_PORT").(string)
	baseRegistryURL = strings.Join([]string{"http://", registryHost, ":", registryPort, "/registry/rest/v1"}, "")

	//check if the URL is accessible
	address := registryHost + ":" + registryPort
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		fmt.Printf("Registry address [%v] is unreachable\n", address)
		os.Exit(1)
	}
}
