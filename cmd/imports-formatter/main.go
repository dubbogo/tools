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

const (
	GO_FILE_SUFFIX = ".go"
	GO_ROOT        = "GOROOT"
	PATH_SEPARATOR = "/"
	IMPORT         = "import"
	GO_MOD         = "go.mod"
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
		fmt.Println("get project name error:", err)
		return
	}

	err = preProcess(goRoot, goPkgMap)
	if err != nil {
		fmt.Println("process go src error:", err)
		return
	}

	err = reformatImports(projectRootPath)
	if err != nil {
		fmt.Println("reformatImports error:", err)
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
		reformatImports(path + "/" + dir.Name())
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

	rootImports := make([]string, 0)
	internalImports := make([]string, 0)
	// import prefix(orgnazation) -> import packages
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

		lineStr := string(line)
		if strings.HasPrefix(lineStr, IMPORT) {
			beginImports = true
		}

		// collect thirdImports
		if beginImports && strings.Contains(lineStr, "\"") {
			importPkg := strings.TrimSpace(lineStr)
			orgImportPkg := importPkg
			// TODO 处理 alias "xxx" 和 "xxx"
			if strings.HasPrefix(importPkg, "_") {
				importPkg = strings.TrimPrefix(importPkg, "_")
				importPkg = strings.TrimSpace(importPkg)
			}
			importPkg = strings.Trim(importPkg, "\"")
			// go root import block
			if _, ok := goPkgMap[importPkg]; ok {
				rootImports = append(rootImports, importPkg)
				continue
			}

			//
			if strings.HasPrefix(importPkg, projectName) {
				internalImports = append(internalImports, importPkg)
				continue
			}

			importsSegment := strings.Split(importPkg, "/")
			project := strings.Join(importsSegment[:3], "/")
			if value, ok := thirdImports[project]; ok {
				value = append(value, orgImportPkg)
				thirdImports[project] = value
			} else {
				thirdImports[project] = []string{orgImportPkg}
			}
		}

		for _, block := range endBlocks {
			if strings.HasPrefix(string(line), block) {
				endImport = true
				beginImports = false
				output = refreshImports(output, rootImports, internalImports, thirdImports)
				break
			}
		}

		if beginImports {
			continue
		}

		output = append(output, line...)
		output = append(output, []byte("\n")...)
	}

	// TODO 覆盖原来文件
	outF, err := os.OpenFile("../test_format.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer outF.Close()
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

func refreshImports(content []byte, rootImports, internalImports []string, thirdImports map[string][]string) []byte {
	if len(rootImports) > 0 {
		content = append(content, []byte("import(\n")...)
		content = doRefreshImports(content, wrapImports(rootImports))
		content = append(content, []byte(")\n\n")...)
	}

	if len(thirdImports) > 0 {
		content = append(content, []byte("import(\n")...)
		thirdProjects := make([]string, 0)
		for key := range thirdImports {
			thirdProjects = append(thirdProjects, key)
		}
		sort.Strings(thirdProjects)
		for idx, key := range thirdProjects {
			value := thirdImports[key]
			// TODO 兼容
			// "xxx" 和 alias "xxx" 这两种格式的排序
			sort.Strings(value)
			content = doRefreshImports(content, value)
			if idx < len(thirdProjects)-1 {
				content = append(content, []byte("\n")...)
			}
		}
		content = append(content, []byte(")\n\n")...)
	}

	if len(internalImports) > 0 {
		content = append(content, []byte("import(\n")...)
		content = doRefreshImports(content, wrapImports(internalImports))
		content = append(content, []byte(")\n\n")...)
	}

	return content
}

func wrapImports(imports []string) []string {
	for i := range imports {
		imports[i] = "\"" + imports[i] + "\""
	}
	return imports
}

func doRefreshImports(content []byte, imports []string) []byte {
	sort.Strings(imports)
	for _, rImport := range imports {
		content = append(content, []byte("\t"+rImport+"\n")...)
	}
	return content
}
