package sample

import (
	"bytes"
	"go/format"
	"os"
	"path"
)

type fileGenerator struct {
	path    string
	file    string
	context string
}

var (
	fileMap = make(map[string]*fileGenerator)
)

func Generate(rootPath string) error {
	for _, v := range fileMap {
		v.path = path.Join(rootPath, v.path)
		if err := genFile(v); err != nil {
			return err
		}
	}
	return nil
}

func genFile(fg *fileGenerator) error {
	fp, err := createFile(fg.path, fg.file)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if _, err := buffer.WriteString(fg.context); err != nil {
		return err
	}
	code := formatCode(buffer.String())
	_, err = fp.WriteString(code)
	return err
}

func createFile(dir, file string) (*os.File, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	return os.Create(path.Join(dir, file))
}

func formatCode(code string) string {
	res, err := format.Source([]byte(code))
	if err != nil {
		return code
	}
	return string(res)
}

func createTemplate(file, content string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	return err
}

func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}
