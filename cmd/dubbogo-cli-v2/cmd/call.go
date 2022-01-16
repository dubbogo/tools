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
	"github.com/dubbogo/tools/internal/client"
	"github.com/dubbogo/tools/internal/json_register"
	"log"

	"github.com/spf13/cobra"
)

// callCmd represents the call command
var (
	callCmd = &cobra.Command{
		Use:   "call",
		Short: "call a method",
		Long:  "",
		Run:   call,
	}
)

var (
	host            string
	port            int
	protocolName    string
	InterfaceID     string
	version         string
	group           string
	method          string
	sendObjFilePath string
	recvObjFilePath string
	timeout         int
)

func init() {
	rootCmd.AddCommand(callCmd)
	showCmd.Flags().String("r", "localhost", "")
	showCmd.Flags().Int("h", 8080, "")

	showCmd.Flags().StringVarP(&host, "h", "", "localhost", "target server host")
	showCmd.Flags().IntVarP(&port, "p", "", 8080, "target server port")
	showCmd.Flags().StringVarP(&protocolName, "proto", "", "dubbo", "transfer protocol")
	showCmd.Flags().StringVarP(&InterfaceID, "i", "", "com", "target service registered interface")
	showCmd.Flags().StringVarP(&version, "v", "", "", "target service version")
	showCmd.Flags().StringVarP(&group, "g", "", "", "target service group")
	showCmd.Flags().StringVarP(&method, "method", "", "", "target method")
	showCmd.Flags().StringVarP(&sendObjFilePath, "sendObj", "", "", "json file path to define transfer struct")
	showCmd.Flags().StringVarP(&recvObjFilePath, "recvObj", "", "", "json file path to define receive struct")
	showCmd.Flags().IntVarP(&timeout, "timeout", "", 3000, "request timeout (ms)")
}

func call(cmd *cobra.Command, args []string) {
	checkParam()
	reqPkg := json_register.RegisterStructFromFile(sendObjFilePath)
	recvPkg := json_register.RegisterStructFromFile(recvObjFilePath)

	t, err := client.NewTelnetClient(host, port, protocolName, InterfaceID, version, group, method, reqPkg, timeout)
	if err != nil {
		panic(err)
	}
	t.ProcessRequests(recvPkg)
	t.Destroy()
}

func checkParam() {
	if method == "" {
		log.Fatalln("-method value not fond")
	}
	if sendObjFilePath == "" {
		log.Fatalln("-sendObj value not found")
	}
	if recvObjFilePath == "" {
		log.Fatalln("-recObj value not found")
	}
}
