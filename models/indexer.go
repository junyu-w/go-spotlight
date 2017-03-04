package models

import (
	"io/ioutil"
	"fmt"
	"path/filepath"
	"os"
)

func IndexAllFiles(dirName string) error {
	absPath, _ := filepath.Abs(dirName)
	curDir, err := ioutil.ReadDir(absPath)
	if err != nil {
		return err
	}
	for _, f := range curDir {
		fmt.Println(f.Name())
		// TODO: get file stat
		if isDir(f) {
			IndexAllFiles(filepath.Join(absPath, f.Name()))
		}
	}
	return nil
}

func isDir(fi os.FileInfo) bool {
	switch mode := fi.Mode(); {
	case mode.IsDir():
		return true
	case mode.IsRegular():
		return false
	default:
		return false
	}
}
