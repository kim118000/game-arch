package file

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
)

func ReadFile(file string) ([]byte, error) {
	b, err := ioutil.ReadFile(file)
	return b, err
}

func FileExt(file string) string {
	return filepath.Ext(file)
}

func FileName(file string) string {
	_, filename := filepath.Split(file)
	return filename
}

func WorkDir() (string, error) {
	path, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		return "", err
	}
	return path, nil
}

func PathJoin(file string) string {
	path, err := WorkDir()
	if err != nil {
		return file
	}
	return filepath.Join(path, file)
}

func ReadDir(dirname string) ([]fs.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	return files, err
}

func ReadDirContent(path string, call func(filename string, content []byte)) error {
	files, err := ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			return ReadDirContent(path+"/"+file.Name(), call)
		} else {
			content, err1 := ioutil.ReadFile(path + "/" + file.Name())
			if err1 != nil {
				return err1
			}
			call(file.Name(), content)
		}
	}
	return nil
}
