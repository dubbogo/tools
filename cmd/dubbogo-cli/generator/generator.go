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
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	hessianImport = []byte{newLine, 105, 109, 112, 111, 114, 116, 32, 40, newLine, 9, 34, 103, 105, 116, 104, 117, 98, 46, 99, 111, 109, 47, 97, 112, 97, 99, 104, 101, 47, 100, 117, 98, 98, 111, 45, 103, 111, 45, 104, 101, 115, 115, 105, 97, 110, 50, 34, newLine, 41, newLine}

	initFunctionPrefix = []byte{newLine, 102, 117, 110, 99, 32, 105, 110, 105, 116, 40, 41, 32, 123, newLine}
	initFunctionSuffix = []byte{newLine, funcEnd, newLine}

	hessianRegistryPOJOFunctionPrefix = []byte{9, 104, 101, 115, 115, 105, 97, 110, 46, 82, 101, 103, 105, 115, 116, 101, 114, 80, 79, 74, 79, 40, 38}
	hessianRegistryPOJOFunctionSuffix = []byte{123, funcEnd, 41}
)

type Generator struct {
	f string
}

func NewGenerator(f string) *Generator {
	return &Generator{f}
}

func (g Generator) Execute() {
	f := g.f
	var file *fileInfo
	var err error
	showLog(infoLog, "=== Generate start [%s] ===", f)
	file, err = scanFile(f)
	if err != nil {
		showLog(errorLog, "=== Generate error [%s] ===\n%v", f, err)
		return
	}
	err = genRegistryStatement(file)
	if err != nil {
		showLog(errorLog, "=== Generate error [%s] ===\n%v", f, err)
		return
	}
	showLog(infoLog, "=== Generate completed [%s] ===", f)
}

// scanFile ??????????????????
func scanFile(filePath string) (file *fileInfo, err error) {
	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	buf := bufio.NewReader(f)
	file = &fileInfo{
		path: filePath,
	}
	bufferSize := 0
	buffer := make([]byte, bufferSize)
	stack := make([][]byte, 0)
	var line []byte
	var lineSize int
	for {
		line, _, err = buf.ReadLine()
		if err == io.EOF {
			err = nil
			break
		}

		bufferSize = len(buffer)
		lineSize = len(line)

		if file.hasInitFunc && lineSize > 0 && line[0] == funcEnd {
			file.initFuncEndIndex = bufferSize + 1 // ??????????????????????????????
		}

		buffer = append(buffer, line...)
		buffer = append(buffer, newLine)

		if passed, _ := regexp.Match(PackageRegexp, line); passed { // ??????package??????
			file.packageStartIndex = bufferSize              // ??????????????????????????????
			file.packageEndIndex = bufferSize + lineSize + 1 // ??????????????????????????????
			continue
		}

		if passed, _ := regexp.Match(InitFunctionRegexp, line); passed { // ?????????????????????
			file.hasInitFunc = true
			file.initFuncStartIndex = bufferSize                     // ??????????????????????????????
			file.initFuncStatementStartIndex = bufferSize + lineSize // ?????????????????????????????????
			continue
		}

		if !file.hasHessianImport {
			r, _ := regexp.Compile(HessianImportRegexp)
			rIndexList := r.FindIndex(line)
			if len(rIndexList) > 0 { // ?????????????????????hessian2???
				checkStatement := line[:rIndexList[0]]
				passed, _ := regexp.Match(LineCommentRegexp, checkStatement) // ??????????????????
				file.hasHessianImport = !passed
				continue
			}
		}

		if passed, _ := regexp.Match(HessianPOJORegexp, line); !passed { // ???????????????Hessian.POJO?????????
			continue
		}
		structName := getStructName(line)
		if len(structName) != 0 {
			stack = append(stack, structName)
		}
	}
	file.buffer = buffer
	file.hessianPOJOList = stack
	return
}

// getStructName ??????Hessian.POJO??????????????????
func getStructName(line []byte) []byte {
	r, _ := regexp.Compile(HessianPOJONameRegexp)
	line = r.Find(line)
	if len(line) != 0 {
		return line[1 : len(line)-1]
	}
	return nil
}

// genRegistryPOJOStatement ??????POJO???????????????
func genRegistryPOJOStatement(pojo []byte) []byte {
	var buffer []byte
	buffer = append(buffer, hessianRegistryPOJOFunctionPrefix...)
	buffer = append(buffer, pojo...)
	buffer = append(buffer, hessianRegistryPOJOFunctionSuffix...)
	return buffer
}

func lastIndexOf(buf []byte, b byte, s *int) int {
	var l int
	if s != nil {
		l = *s
	} else {
		l = len(buf)
	}
	if l < 0 || l >= len(buf) {
		return -1
	}
	for ; l >= 0; l-- {
		if buf[l] == b {
			return l
		}
	}
	return -1
}

func escape(str string) string {
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, ".", "\\.")
	str = strings.ReplaceAll(str, "(", "\\(")
	str = strings.ReplaceAll(str, ")", "\\)")
	str = strings.ReplaceAll(str, "{", "\\{")
	str = strings.ReplaceAll(str, "}", "\\}")
	return str
}

// genRegistryPOJOStatements ??????POJO???????????????
func genRegistryPOJOStatements(file *fileInfo, initFunctionStatement *[]byte) []byte {
	f := file.path
	hessianPOJOList := file.hessianPOJOList
	var buffer []byte
	var r *regexp.Regexp
	var rIndexList []int
	for _, name := range hessianPOJOList {
		statement := genRegistryPOJOStatement(name)
		if initFunctionStatement != nil { // ????????????????????????????????????
			r, _ = regexp.Compile(escape(string(statement)))
			initStatement := *initFunctionStatement
			rIndexList = r.FindIndex(initStatement)
			if len(rIndexList) > 0 {
				i := rIndexList[0]
				n := lastIndexOf(initStatement, newLine, &i)
				checkStatement := initStatement[lastIndexOf(initStatement, newLine, &n)+1 : i]
				if passed, _ := regexp.Match(LineCommentRegexp, checkStatement); !passed { // ??????????????????
					showLog(infoLog, "=== Ignore POJO [%s].%s ===", f, name)
					continue // ???????????????????????????
				}
			}
		}

		showLog(infoLog, "=== Registry POJO [%s].%s ===", f, name)

		buffer = append(buffer, newLine)
		buffer = append(buffer, statement...)
		buffer = append(buffer, newLine)
	}
	return buffer
}

func genRegistryStatement(file *fileInfo) error {
	if file == nil || len(file.hessianPOJOList) == 0 {
		return nil
	}
	buffer := file.buffer

	offset := 0

	if !file.hasHessianImport { // ??????hessian2??????
		sliceIndex := file.packageEndIndex + offset
		var bufferClone []byte
		bufferClone = append(bufferClone, buffer[:sliceIndex]...)
		bufferClone = append(bufferClone, hessianImport...)
		bufferClone = append(bufferClone, buffer[sliceIndex:]...)
		buffer = bufferClone
		offset += len(hessianImport)
	}

	var registryPOJOStatement []byte

	if file.hasInitFunc {
		sliceIndex := file.initFuncStatementStartIndex + offset
		initFunctionStatement := buffer[sliceIndex : file.initFuncEndIndex+offset-1]
		registryPOJOStatement = genRegistryPOJOStatements(file, &initFunctionStatement)
		var bufferClone []byte
		bufferClone = append(bufferClone, buffer[:sliceIndex]...)
		bufferClone = append(bufferClone, registryPOJOStatement...) // ??????POJO??????????????????init??????
		bufferClone = append(bufferClone, buffer[sliceIndex:]...)
		buffer = bufferClone
	} else { // ?????????????????????
		registryPOJOStatement = genRegistryPOJOStatements(file, nil)
		buffer = append(buffer, initFunctionPrefix...)    // ??????init??????
		buffer = append(buffer, registryPOJOStatement...) // ??????POJO??????????????????init??????
		buffer = append(buffer, initFunctionSuffix...)
	}

	offset += len(registryPOJOStatement)

	f, err := os.OpenFile(file.path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	_, err = f.Write(buffer)
	return err
}
