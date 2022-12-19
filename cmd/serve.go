/*
 * Copyright (c) 2022 Michaël COLL.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"strconv"

	"github.com/michaelcoll/weasyprintaas/internal/weasyprint"
)

const (
	portEnvVarName = "WPR_PORT"
)

var port uint16
var multithreading bool

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		weasyprint.New(multithreading).Serve(getPort())
	},
}

func init() {
	serveCmd.Flags().Uint16VarP(&port, "port", "p", 8080, "Listening port")
	serveCmd.Flags().BoolVarP(&multithreading, "multi-threading", "T", false, "Multi-threading for conversions")

	rootCmd.AddCommand(serveCmd)
}

func getPort() uint16 {
	env, present := os.LookupEnv(portEnvVarName)
	if present {
		if p, err := strconv.ParseUint(env, 10, 16); err == nil {
			return uint16(p)
		}
	}

	return port
}
