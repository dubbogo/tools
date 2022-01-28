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

package telnet

import (
	"log"
)

import (
	"github.com/dubbogo/tools/cmd/dubbogo-cli/common"
	"github.com/dubbogo/tools/internal/client"
	"github.com/dubbogo/tools/internal/json_register"
)

type TelnetAdapter struct {
	Host            string
	Port            int
	ProtocolName    string
	InterfaceID     string
	Version         string
	Group           string
	Method          string
	SendObjFilePath string
	RecvObjFilePath string
	Timeout         int
}

func (a *TelnetAdapter) CheckParam() bool {
	if a.Method == "" {
		log.Fatalln("-method value not fond")
	}
	if a.SendObjFilePath == "" {
		log.Fatalln("-sendObj value not found")
	}
	if a.RecvObjFilePath == "" {
		log.Fatalln("-recObj value not found")
	}
	return true
}

func (a *TelnetAdapter) Execute() {
	reqPkg := json_register.RegisterStructFromFile(a.SendObjFilePath)
	recvPkg := json_register.RegisterStructFromFile(a.RecvObjFilePath)

	t, err := client.NewTelnetClient(a.Host, a.Port, a.ProtocolName, a.InterfaceID, a.Version, a.Group, a.Method, reqPkg, a.Timeout)
	if err != nil {
		panic(err)
	}
	t.ProcessRequests(recvPkg)
	t.Destroy()
}

func (a *TelnetAdapter) GetMode() common.AdapterMode {
	return common.TelnetAdapterMode
}
