package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Video struct {
	name, dir, path string
	start           uint
}

const (
	DIR_SEPARATOR byte = ':'
)

func (v *Video) Find() error {

	if strings.IndexByte(v.dir, DIR_SEPARATOR) >= 0 {
		v.dir = strings.ReplaceAll(v.dir, string(DIR_SEPARATOR), "/")
	}

	files, err := GetFiles(v.dir)
	if err != nil {
		return err
	}

	var fullName string

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if strings.Index(file.Name(), v.name) == 0 {
			fullName = file.Name()
		}

	}

	if len(fullName) == 0 {
		return fmt.Errorf("file not found: %s", v.name)
	}

	v.path = filepath.Join(workingDir, v.dir, fullName)

	return nil

}

func GetFiles(dir string) ([]os.DirEntry, error) {

	files := dirCache[dir]
	if len(files) != 0 {
		return files, nil
	}

	searchDir := workingDir + dir
	files, err := os.ReadDir(searchDir)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("empty directory: %s", dir)
	}

	dirCache[dir] = files
	return files, nil

}
