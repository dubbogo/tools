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

package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/dubbogo/tools/cmd/dubbogo-cli/common"
)

type GeneratorAdapter struct {
	DirPath     string // 文件扫描目录
	ThreadLimit int    // 最大工作线程数
}

func (a *GeneratorAdapter) CheckParam() bool {
	if a.ThreadLimit < 1 {
		a.ThreadLimit = 1
	}
	return true
}

func (a *GeneratorAdapter) Execute() {
	fileList, err := listFiles(a.DirPath, targetFileSuffix)
	if err != nil {
		showLog(errorLog, "%v", err)
		return
	}

	pool := NewPool(a.ThreadLimit)
	for _, f := range fileList {
		pool.Execute(NewGenerator(f))
	}
	pool.Wait()
}

func (a *GeneratorAdapter) GetMode() common.AdapterMode {
	return common.GeneratorAdapterMode
}

type fileInfo struct {
	path string

	packageStartIndex int
	packageEndIndex   int

	hasInitFunc                 bool
	initFuncStartIndex          int
	initFuncEndIndex            int
	initFuncStatementStartIndex int

	hasHessianImport bool

	buffer []byte

	hessianPOJOList [][]byte
}

// listFiles 获取目标目录下所有go文件
func listFiles(dirPath string, suffix string) (fileList []string, err error) {
	suffix = strings.ToUpper(suffix)
	fileList = make([]string, 0)
	_, err = os.Lstat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("找不到该目录[%s]", dirPath)
	}
	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error { // 递归获取目录下所有go文件
		if d == nil || d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(d.Name()), suffix) {
			fileList = append(fileList, path)
		}
		return nil
	})
	return
}
