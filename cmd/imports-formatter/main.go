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
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/tools/constant"
)

const (
	GO_FILE_SUFFIX = ".go"
	GO_ROOT        = "GOROOT"
	QUOTATION_MARK = "\""
	PATH_SEPARATOR = "/"
	IMPORT         = "import"
	GO_MOD         = "go.mod"
)

var (
	blankLine         bool
	currentWorkDir, _ = os.Getwd()
	goRoot            = os.Getenv(GO_ROOT) + "/src"
	endBlocks         = []string{"var", "const", "type", "func"}
	projectRootPath   string
	projectName       string
	goPkgMap          = make(map[string]struct{})
	outerComments     = make([]string, 0)
	// record comments between importBlocks and endBlocks
	innerComments = make([]string, 0)
	ignorePath    = []string{".git", ".idea", ".github", ".vscode", "vendor", "swagger", "docs"}
	newLine       = false
	blockCount    = 0
)

func init() {
	flag.StringVar(&projectRootPath, "path", currentWorkDir, "the path need to be reformatted")
	flag.StringVar(&projectName, "module", "", "project name, namely module name in the go.mod")
	flag.BoolVar(&blankLine, "bl", true, "if true, it will split different import modules with a blank line")
}

func main() {
	fmt.Println("imports-formatter:", constant.Version)
	flag.Parse()
	var err error
	projectName, err = getProjectName(projectRootPath)
	if err != nil {
		panic(err)
	}

	err = preProcess(goRoot, goPkgMap)
	if err != nil {
		panic(err)
	}

	err = reformatImports(projectRootPath)
	if err != nil {
		panic(err)
	}
}

func getProjectName(path string) (string, error) {
	if projectName != "" {
		return projectName, nil
	}

	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == GO_MOD {
			f, err := os.OpenFile(path+PATH_SEPARATOR+fileInfo.Name(), os.O_RDONLY, 0644)
			if err != nil {
				return "", err
			}

			reader := bufio.NewReader(f)
			for {
				line, _, err := reader.ReadLine()
				if err != nil {
					return "", err
				}
				lineStr := strings.TrimSpace(string(line))
				if strings.HasPrefix(lineStr, "module") {
					return strings.Split(lineStr, " ")[1], nil
				}
			}
		}
	}

	return "", err
}

func preProcess(path string, goPkgMap map[string]struct{}) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.WithStack(err)
	}

	dirs := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo)
		} else if strings.HasSuffix(fileInfo.Name(), GO_FILE_SUFFIX) {
			goPkgMap[strings.TrimPrefix(path, goRoot+"/")] = struct{}{}
		}
	}

	for _, dir := range dirs {
		err := preProcess(path+PATH_SEPARATOR+dir.Name(), goPkgMap)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func reformatImports(path string) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return errors.WithStack(err)
	}

	dirs := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() && !ignore(fileInfo.Name()) {
			dirs = append(dirs, fileInfo)
		} else if strings.HasSuffix(fileInfo.Name(), GO_FILE_SUFFIX) {
			clearData()
			newLine = false
			err = doReformat(path + PATH_SEPARATOR + fileInfo.Name())
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	for _, dir := range dirs {
		err := reformatImports(path + "/" + dir.Name())
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func doReformat(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(errors.New(filePath + "encounter error:" + err.Error()))
		}
	}(f)
	if err != nil {
		return errors.New("open " + filePath + " encounter error:" + err.Error())
	}

	reader := bufio.NewReader(f)
	beginImports := false
	endImport := false
	output := make([]byte, 0)

	// processed import(organization) -> original import packages
	rootImports := make(map[string][]string)
	internalImports := make(map[string][]string)
	thirdImports := make(map[string][]string)

	for {
		if len(outerComments) > 0 {
			for _, c := range outerComments {
				output = append(output, []byte(c+"\n")...)
			}
			outerComments = make([]string, 0)
		}

		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				root := len(rootImports)
				internal := len(internalImports)
				third := len(thirdImports)
				if root > 0 && internal > 0 && third > 0 {
					blockCount = 3
				} else if (root > 0 && internal > 0) || (root > 0 && third > 0) || (internal > 0 && third > 0) {
					blockCount = 2
				} else if root > 0 || internal > 0 || third > 0 {
					blockCount = 1
				}
				break
			}
			return errors.New("read line of " + filePath + " encounter error:" + err.Error())
		}

		if endImport {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
			continue
		}

		// if import blocks end
		for _, block := range endBlocks {
			if strings.HasPrefix(string(line), block) {
				endImport = true
				beginImports = false
				newLine = true
				output = refreshImports(output, mergeImports(rootImports), false)
				output = refreshImports(output, mergeImports(thirdImports), blankLine)
				output = refreshImports(output, mergeImports(internalImports), false)
				if len(innerComments) > 0 {
					for _, c := range innerComments {
						output = append(output, []byte(c+"\n")...)
					}
				}
				break
			}
		}

		lineStr := string(line)
		if strings.HasPrefix(lineStr, IMPORT) {
			beginImports = true
		}

		orgImportPkg := strings.TrimSpace(lineStr)
		// single line comment
		if strings.HasPrefix(orgImportPkg, "//") || (strings.HasPrefix(orgImportPkg, "/*") && strings.HasSuffix(orgImportPkg, "*/")) {
			if beginImports {
				innerComments = append(innerComments, lineStr)
			} else {
				outerComments = append(outerComments, lineStr)
			}
			continue
		}
		// multiple lines comment
		if strings.HasPrefix(orgImportPkg, "/*") {
			if beginImports {
				innerComments = append(innerComments, lineStr)
				commentLine, _, err := reader.ReadLine()
				commentLineStr := string(commentLine)
				for err == nil && !strings.HasSuffix(strings.TrimSpace(commentLineStr), "*/") {
					innerComments = append(innerComments, commentLineStr)
					commentLine, _, err = reader.ReadLine()
					commentLineStr = string(commentLine)
				}
				if err == nil {
					innerComments = append(innerComments, commentLineStr)
				} else {
					return errors.New("read line of " + filePath + " encounter error:" + err.Error())
				}
			} else {
				outerComments = append(outerComments, lineStr)
				commentLine, _, err := reader.ReadLine()
				commentLineStr := string(commentLine)
				for err == nil && !strings.HasSuffix(strings.TrimSpace(commentLineStr), "*/") {
					outerComments = append(outerComments, commentLineStr)
					commentLine, _, err = reader.ReadLine()
					commentLineStr = string(commentLine)
				}
				if err == nil {
					outerComments = append(outerComments, commentLineStr)
				} else {
					return errors.New("read line of " + filePath + " encounter error:" + err.Error())
				}
			}
			continue
		}

		// collect imports
		if beginImports && strings.Contains(orgImportPkg, QUOTATION_MARK) {
			innerComments = innerComments[:0]
			// single line import
			if strings.HasPrefix(orgImportPkg, IMPORT+" ") {
				orgImportPkg = strings.TrimPrefix(orgImportPkg, IMPORT+" ")
			}
			importKey := orgImportPkg
			// process those imports that has alias
			importKey = unwrapImport(importKey)

			if _, ok := goPkgMap[importKey]; ok {
				// go root import block
				cacheImports(rootImports, importKey, []string{orgImportPkg})
			} else if importKey == projectName || (strings.HasPrefix(importKey, projectName) && len(importKey) > len(projectName) && importKey[len(projectName)] == '/') {
				/**
				for project a
				****************************
				import (
					a
					a/b
					aa
				)
				****************************
				we need to combine a&a/b, and exclude aa.
				importKey == projectName is for a
				strings.HasPrefix(importKey, projectName) && len(importKey) > len(projectName) && importKey[len(projectName)] == '/' is for a/b
				if we simply use strings.HasPrefix(importKey, projectName), it will recognize a&aa as the same project.
				*/
				// internal imports of the project
				cacheImports(internalImports, importKey, []string{orgImportPkg})
			} else {
				// imports of the third projects
				project, importsSegment := "", strings.Split(importKey, "/")

				// like google.golang.org/grpc etc.
				if len(importsSegment) == 2 {
					project = strings.Join(importsSegment[:2], "/")
				} else if len(importsSegment) > 2 {
					project = strings.Join(importsSegment[:3], "/")
				} else {
					return errors.New("unexpected import format: " + orgImportPkg + " in file " + filePath)
				}
				cacheImports(thirdImports, project, []string{orgImportPkg})
			}
			continue
		}

		// to process `import (`
		if beginImports {
			continue
		}

		output = append(output, line...)
		output = append(output, []byte("\n")...)
	}

	if !endImport {
		output = refreshImports(output, mergeImports(rootImports), false)
		output = refreshImports(output, mergeImports(thirdImports), blankLine)
		output = refreshImports(output, mergeImports(internalImports), false)
		if len(innerComments) > 0 {
			for _, c := range innerComments {
				output = append(output, []byte(c+"\n")...)
			}
		}
	}

	outF, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.New("open/create" + filePath + " encounter error:" + err.Error())
	}
	defer func(outF *os.File) {
		err := outF.Close()
		if err != nil {

		}
	}(outF)
	writer := bufio.NewWriter(outF)
	_, err = writer.Write(output)
	if err != nil {
		return errors.New("write " + filePath + " encounter error:" + err.Error())
	}
	err = writer.Flush()
	if err != nil {
		return errors.New("flush " + filePath + " encounter error:" + err.Error())
	}
	return nil
}

func unwrapImport(importStr string) string {
	// exists alias
	if !strings.HasPrefix(importStr, QUOTATION_MARK) {
		importStr = strings.Fields(importStr)[1]
	}
	return strings.Trim(importStr, QUOTATION_MARK)
}

func cacheImports(m map[string][]string, key string, values []string) {
	if oldValues, ok := m[key]; ok {
		oldValues = append(oldValues, values...)
		m[key] = oldValues
	} else {
		m[key] = values
	}
}

func mergeImports(m map[string][]string) map[string][]string {
	mergedMap := make(map[string][]string)
	for key := range m {
		mergedKeys := make([]string, len(m))
		newMergedMap := make(map[string][]string)
		mergedValues := make([]string, 0)
		mergedValues = append(mergedValues, m[key]...)
		rootKey := key

		for mKey := range mergedMap {
			if strings.HasPrefix(rootKey, mKey) {
				rootKey = mKey
			}
			if strings.HasPrefix(key, mKey) || strings.HasPrefix(mKey, key) {
				// mKey is a sub package of the module key || key is a sub package of the module mKey
				mergedKeys = append(mergedKeys, mKey)
			}
		}

		for mKey, value := range mergedMap {
			target := false
			for _, mKey1 := range mergedKeys {
				if mKey == mKey1 {
					target = true
					mergedValues = append(mergedValues, mergedMap[mKey]...)
					break
				}
			}

			if !target {
				cacheImports(newMergedMap, mKey, value)
			}
		}

		cacheImports(newMergedMap, rootKey, mergedValues)
		mergedMap = newMergedMap
	}

	return mergedMap
}

func refreshImports(content []byte, importsMap map[string][]string, blankLine bool) []byte {
	if len(importsMap) <= 0 {
		return content
	}
	blockCount--
	content = append(content, []byte("import (\n")...)
	sortedKeys := make([]string, 0, len(importsMap))
	for key := range importsMap {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	for idx, key := range sortedKeys {
		value := importsMap[key]
		content = doRefreshImports(content, value)
		if blankLine && idx < len(sortedKeys)-1 {
			content = append(content, []byte("\n")...)
		}
	}
	if !newLine && blockCount == 0 && len(innerComments) <= 0 {
		content = append(content, []byte(")\n")...)
	} else {
		content = append(content, []byte(")\n\n")...)
	}
	return content
}

func doRefreshImports(content []byte, imports []string) []byte {
	sort.SliceStable(imports, func(i, j int) bool {
		v1 := unwrapImport(imports[i])
		v2 := unwrapImport(imports[j])
		return v1 < v2
	})
	for _, rImport := range imports {
		content = append(content, []byte("\t"+rImport+"\n")...)
	}
	return content
}

func clearData() {
	innerComments = innerComments[:0]
	outerComments = outerComments[:0]
}

func ignore(path string) bool {
	for _, name := range ignorePath {
		if name == path {
			return true
		}
	}
	return false
}
