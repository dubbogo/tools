package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	include   = flag.String("include", "./", "hessian映射关系目录，默认为项目根目录")
	onlyError = flag.Bool("error", false, "只输出错误信息，默认全量输出")
	maxThread = flag.Int("thread", runtime.NumCPU()*2, "最大工作线程数，默认为CPU核数的2倍")
)

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

//go:generate go fmt
func main() {
	flag.Parse()

	fileList, err := listFiles(*include, targetFileSuffix)
	if err != nil {
		showLog(errorLog, "%v", err)
		return
	}

	pool := NewPool(*maxThread)
	for _, f := range fileList {
		pool.Execute(NewGenerator(f))
	}
	pool.Wait()
}
