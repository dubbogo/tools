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

package main

import (
	"flag"
	"runtime"
)

import (
	"github.com/dubbogo/tools/cmd/dubbogo-cli/common"
	"github.com/dubbogo/tools/cmd/dubbogo-cli/generator"
	"github.com/dubbogo/tools/cmd/dubbogo-cli/telnet"
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

	isRegistryGeneratorMode    bool   // 是否为生成器模式
	registryFileDirectoryPath  string // 文件扫描目录
	generatorWorkerThreadLimit int    // 最大工作线程数
)

var (
	telnetAdapter    common.Adapter
	generatorAdapter common.Adapter

	adapters []common.Adapter
)

func init() {
	flag.StringVar(&host, "h", "localhost", "target server host")
	flag.IntVar(&port, "p", 8080, "target server port")
	flag.StringVar(&protocolName, "proto", "dubbo", "transfer protocol")
	flag.StringVar(&InterfaceID, "i", "com", "target service registered interface")
	flag.StringVar(&version, "v", "", "target service version")
	flag.StringVar(&group, "g", "", "target service group")
	flag.StringVar(&method, "method", "", "target method")
	flag.StringVar(&sendObjFilePath, "sendObj", "", "json file path to define transfer struct")
	flag.StringVar(&recvObjFilePath, "recvObj", "", "json file path to define receive struct")
	flag.IntVar(&timeout, "timeout", 3000, "request timeout (ms)")

	telnetAdapter = &telnet.TelnetAdapter{
		Host:            host,
		Port:            port,
		ProtocolName:    protocolName,
		InterfaceID:     InterfaceID,
		Version:         version,
		Group:           group,
		Method:          method,
		SendObjFilePath: sendObjFilePath,
		RecvObjFilePath: recvObjFilePath,
		Timeout:         timeout,
	}
	adapters = append(adapters, telnetAdapter)

	flag.BoolVar(&isRegistryGeneratorMode, "generator", false, "switch to registry statement generator mode, default `false`")
	flag.StringVar(&registryFileDirectoryPath, "include", "./", "file scan directory path, default `./`")
	flag.IntVar(&generatorWorkerThreadLimit, "thread", runtime.NumCPU()*2, "worker thread limit, default (cpu core) * 2")

	generatorAdapter = &generator.GeneratorAdapter{
		DirPath:     registryFileDirectoryPath,
		ThreadLimit: generatorWorkerThreadLimit,
	}
	adapters = append(adapters, generatorAdapter)
}

func chooseAdapter(mode common.AdapterMode) common.Adapter {
	for _, adapter := range adapters {
		if adapter.GetMode() == mode {
			return adapter
		}
	}
	return telnetAdapter
}

func main() {
	flag.Parse()

	var mode common.AdapterMode
	if isRegistryGeneratorMode {
		mode = common.GeneratorAdapterMode
	}
	adapter := chooseAdapter(mode)
	if adapter.CheckParam() {
		adapter.Execute()
	}
}
