package models

import (
	"os"
	"path/filepath"
	"time"
)

const (
	initReadLength int64 = 500 // read last 500 bytes
)

type FileRecord struct {
	Path              string
	Name              string
	Extension         string
	AccessTime        time.Time
	ModifyTime        time.Time // content was modified
	ChangeTime        time.Time
	RecentWords       string
	LastIndexPosition int64
}

func (fr *FileRecord) ReadAt(p []byte, off int64) (int, error) {
	f, err := os.Open(fr.Path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.ReadAt(p, off)
}

func NewFileRecord(absPath string, fName string, atime, mtime, ctime time.Time, size int64) *FileRecord {
	lastIndexPosition := size - (initReadLength % size)
	fr := &FileRecord{
		Path:              absPath,
		Name:              fName,
		Extension:         filepath.Ext(fName),
		AccessTime:        atime,
		ModifyTime:        mtime,
		ChangeTime:        ctime,
		LastIndexPosition: lastIndexPosition,
	}
	fr.initRecentWords()
	return fr
}

func (fr *FileRecord) initRecentWords() {
	buffer := make([]byte, initReadLength, initReadLength)
	l, err := fr.ReadAt(buffer, fr.LastIndexPosition)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	fr.RecentWords = string(buffer)
	fr.LastIndexPosition += int64(l)
}
