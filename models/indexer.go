package models

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	indexFileError string = "ERROR: failed to index file %s"
)

// ##################### Manual Indexing ########################

func statTimes(fi os.FileInfo) (atime, mtime, ctime time.Time, err error) {
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Stat_t)
	atime = time.Unix(int64(stat.Atimespec.Sec), int64(stat.Atimespec.Nsec))
	ctime = time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec))
	return
}

func indexFile(fi os.FileInfo, absPath string, fr_idx *bleve.Index) {
	atime, mtime, ctime, err := statTimes(fi)
	if err != nil {
		panic(fmt.Errorf(indexFileError, fi.Name()))

	}
	fr := NewFileRecord(absPath, fi.Name(), atime, mtime, ctime, fi.Size())
	// TODO: use batch index
	(*fr_idx).Index(fr.Path, fr)
}

func IndexAllFiles(dirName string, fr_index *bleve.Index) error {
	absPath, _ := filepath.Abs(dirName)
	curDir, err := ioutil.ReadDir(absPath)
	if err != nil {
		return err
	}
	for _, fi := range curDir {
		// TODO: aysnc index with gorouting
		absFilePath := filepath.Join(absPath, fi.Name())
		if isDir(fi) {
			IndexAllFiles(absFilePath, fr_index)
		} else {
			indexFile(fi, absFilePath, fr_index)
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
		return true
	}
}
