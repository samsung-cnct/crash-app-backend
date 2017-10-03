// Copyright © 2017 Samsung CNCT
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

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "crashbackend",
	Short: "crash app backend validating proxy",
	Long:  `The crash app backend is a reverse proxy that validates crash app requests to elasticsearch`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//backendproxy.Verbose = Verbose
	},

	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	//RootCmd.PersistentFlags().StringVarP(&Target, "target", "t", "", "If set the reverse proxy target will be set to this.")
	//RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "More verbose output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
