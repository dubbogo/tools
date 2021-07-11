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
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
)

const (
	GO_FILE_SUFFIX  = ".go"
	GO_ROOT         = "GOROOT"
	ALIAS_SEPARATOR = " "
	PATH_SEPARATOR  = "/"
	IMPORT          = "import"
	GO_MOD          = "go.mod"
)

var (
	currentWorkDir, _ = os.Getwd()
	goRoot            = os.Getenv(GO_ROOT) + "/src"
	endBlocks         = []string{"var", "const", "type", "func"}
	projectRootPath   string
	projectName       string
	goPkgMap          = make(map[string]struct{})
)

func init() {
	flag.StringVar(&projectRootPath, "path", currentWorkDir, "the path need to be reformatted")
	flag.StringVar(&projectName, "module", "", "project name, namely module name in the go.mod")
}

func main() {
	flag.Parse()
	var err error
	projectName, err = getProjectName(projectRootPath)
	if err != nil {
		panic(err)
		return
	}

	err = preProcess(goRoot, goPkgMap)
	if err != nil {
		panic(err)
		return
	}

	err = reformatImports(projectRootPath)
	if err != nil {
		panic(err)
		return
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
		return err
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
			return err
		}
	}

	return nil
}

func reformatImports(path string) error {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	dirs := make([]os.FileInfo, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirs = append(dirs, fileInfo)
		} else if strings.HasSuffix(fileInfo.Name(), GO_FILE_SUFFIX) {
			err = doReformat(path + PATH_SEPARATOR + fileInfo.Name())
			if err != nil {
				return err
			}
		}
	}

	for _, dir := range dirs {
		err := reformatImports(path + "/" + dir.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func doReformat(filePath string) error {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	beginImports := false
	endImport := false
	output := make([]byte, 0)

	// processed import(orgnazation) -> orignal import packages
	rootImports := make(map[string][]string)
	internalImports := make(map[string][]string)
	thirdImports := make(map[string][]string)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if endImport {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
			continue
		}

		// if imports blocks end
		for _, block := range endBlocks {
			if strings.HasPrefix(string(line), block) {
				endImport = true
				beginImports = false
				output = refreshImports(output, mergeImports(rootImports), false)
				output = refreshImports(output, mergeImports(thirdImports), true)
				output = refreshImports(output, mergeImports(internalImports), false)
				break
			}
		}

		lineStr := string(line)
		if strings.HasPrefix(lineStr, IMPORT) {
			beginImports = true
		}

		// collect imports
		if beginImports && strings.Contains(lineStr, "\"") {
			orgImportPkg := strings.TrimSpace(lineStr)
			if strings.HasPrefix(orgImportPkg, "//") {
				continue
			}
			if strings.HasPrefix(orgImportPkg, "import ") {
				orgImportPkg = strings.TrimPrefix(orgImportPkg, "import ")
			}
			importKey := orgImportPkg
			// process those imports that has alias
			importKey = unwrapImport(importKey)

			if _, ok := goPkgMap[importKey]; ok {
				// go root import block
				cacheImports(rootImports, importKey, []string{orgImportPkg})
			} else if strings.HasPrefix(importKey, projectName) {
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

	outF, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func(outF *os.File) {
		err := outF.Close()
		if err != nil {

		}
	}(outF)
	writer := bufio.NewWriter(outF)
	_, err = writer.Write(output)
	if err != nil {
		return err
	}
	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}

func unwrapImport(importStr string) string {
	if strings.Index(importStr, ALIAS_SEPARATOR) != -1 {
		importStr = strings.Fields(importStr)[1]
	}
	return strings.Trim(importStr, "\"")
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
		merged := false
		for mergedKey := range mergedMap {
			if strings.HasPrefix(key, mergedKey) {
				// key is a sub package of the module mergedKey
				cacheImports(mergedMap, mergedKey, m[key])
				merged = true
			} else if strings.HasPrefix(mergedKey, key) {
				// mergedKey is a sub package of the module key
				mergedValues := mergedMap[mergedKey]
				delete(m, mergedKey)
				cacheImports(mergedMap, key, append(m[key], mergedValues...))
				merged = true
			}
		}
		if merged {
			continue
		}
		cacheImports(mergedMap, key, m[key])
	}

	return mergedMap
}

func refreshImports(content []byte, importsMap map[string][]string, blankLine bool) []byte {
	if len(importsMap) <= 0 {
		return content
	}

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

	content = append(content, []byte(")\n\n")...)
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
