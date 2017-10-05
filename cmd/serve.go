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

	"github.com/oneilcin/crash-app-backend/backendproxy"
	"github.com/spf13/cobra"
)

// Target for the reverse proxy
var Target string
var klogMax int
var ktaskMax int
var rateLimitPerMin int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "starts the backend proxy server",
	Long: `The serve command starts the backend proxy server.  You are required to pass in a target argument.
Example Usage:
crashbackend serve --target "http://elasticsearch"
crashbackend serve --target "http://elasticsearch" --logmax 2000000 --taskmax 300

`,
	Run: func(cmd *cobra.Command, args []string) {
		if Target != "" {
			// must enter a target arg
			fmt.Println(Target)
			backendproxy.Server(Target, klogMax, ktaskMax, rateLimitPerMin)
		} else {
			cmd.Usage()
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")
	serveCmd.PersistentFlags().StringVarP(&Target, "target", "t", "", "The target reverse proxy URL will be set to this.")
	serveCmd.Flags().IntVarP(&klogMax, "logmax", "", 2000000, "Maximum characters for validating KrakenLog string size.")
	serveCmd.Flags().IntVarP(&ktaskMax, "taskmax", "", 400, "Maximum characters for validating KrakenTask string size.")
	serveCmd.Flags().IntVarP(&rateLimitPerMin, "ratelimit", "", 60, "Rate Limit Per Minute Per Handler")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//	serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//RootCmd.Flags().StringVarP(&Source, "source", "s", "", "Source directory to read from")
}
