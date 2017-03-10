package models

import (
	"os"
	"time"
)

const (
	initReadLength int64 = 500 // read last 500 bytes
)

type FileRecord struct {
	Path              string
	Name              string
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
	return f.ReadAt(p, off)
}

func NewFileRecord(absPath string, fName string, atime, mtime, ctime time.Time) *FileRecord {
	fi, _ := os.Stat(absPath)
	lastIndexPosition := fi.Size() - (initReadLength % fi.Size())
	fr := &FileRecord{
		Path:              absPath,
		Name:              fName,
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
