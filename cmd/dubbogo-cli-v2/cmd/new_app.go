/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
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
)

import (
	"github.com/dubbogo/tools/cmd/dubbogo-cli-v2/generator/application"
)

// newAppCmd represents the new command
var newAppCmd = &cobra.Command{
	Use:   "newApp",
	Short: "new a dubbo-go application project",
	Run:   createApp,
}

func init() {
	rootCmd.AddCommand(newAppCmd)
	newAppCmd.Flags().String("path", "rootPath", "")
}

func createApp(cmd *cobra.Command, _ []string) {
	path, err := cmd.Flags().GetString("path")
	if err != nil {
		panic(err)
	}
	if err := application.Generate(path); err != nil {
		panic(err)
	}
}
