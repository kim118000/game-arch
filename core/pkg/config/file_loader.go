package config

import (
	"github.com/kim118000/core/toolkit/file"
)

type FileLoader struct {
	pathRoot string
}

func NewFileLoader(path string) *FileLoader {
	return &FileLoader{
		pathRoot: path,
	}
}

func (f *FileLoader) Load(name string) ([]byte, error) {
	b, err := file.ReadFile(f.pathRoot + "/" + name)
	if err != nil {
		return nil, err
	}
	return b, nil
}