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

const BATCH_SIZE int = 400

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

func StartIndexing(dirName string, fr_index bleve.Index, doneChan chan bool) {
	quitChan := make(chan bool)
	frChan := make(chan *FileRecord)
	go func(fr_index bleve.Index, frChan chan *FileRecord, quitChan chan bool, doneChan chan bool) {
		batch := fr_index.NewBatch()
		batchCount := 0
		for {
			select {
			case fr := <-frChan:
				batch.Index(fr.Path, fr)
				batchCount++
				if batchCount >= BATCH_SIZE {
					fmt.Println("Indexed ", BATCH_SIZE, " files")
					err := fr_index.Batch(batch)
					if err != nil {
						panic(err)
					}
					batch = fr_index.NewBatch()
					batchCount = 0
				}
			case quitSignal := <-quitChan:
				if quitSignal {
					fmt.Println("Indexed ", batch.Size(), " files")
					// index last batch then quit
					err := fr_index.Batch(batch)
					if err != nil {
						panic(err)
					}
					doneChan <- true
					return
				}
			}
		}
	}(fr_index, frChan, quitChan, doneChan)
	IndexAllFiles(dirName, frChan)
	quitChan <- true
}

func enqueFile(fi os.FileInfo, absPath string, bufferChan chan *FileRecord) {
	atime, mtime, ctime, err := statTimes(fi)
	if err != nil {
		panic(fmt.Errorf(indexFileError, fi.Name()))
	}
	fr := NewFileRecord(absPath, fi.Name(), atime, mtime, ctime, fi.Size())
	bufferChan <- fr
}

func IndexAllFiles(dirName string, bufferChan chan *FileRecord) error {
	absPath, _ := filepath.Abs(dirName)
	curDir, err := ioutil.ReadDir(absPath)
	if err != nil {
		return err
	}
	for _, fi := range curDir {
		// TODO: aysnc index with gorouting
		absFilePath := filepath.Join(absPath, fi.Name())
		if isDir(fi) {
			IndexAllFiles(absFilePath, bufferChan)
		} else {
			enqueFile(fi, absFilePath, bufferChan)
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
