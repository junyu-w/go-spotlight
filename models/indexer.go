package models

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/cheggaaa/pb"
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

	progBar := pb.New(getNumFiles(dirName))
	progBar.ShowCounters = true
	progBar.SetWidth(80)
	progBar.Start()

	go func(fr_index bleve.Index) {
		batch := fr_index.NewBatch()
		batchCount := 0
		for {
			select {
			case fr := <-frChan:
				progBar.Increment()
				batch.Index(fr.Path, fr)
				batchCount++
				if batchCount >= BATCH_SIZE {
					err := fr_index.Batch(batch)
					if err != nil {
						panic(err)
					}
					batch = fr_index.NewBatch()
					batchCount = 0
				}
			case quitSignal := <-quitChan:
				if quitSignal {
					// index last batch then quit
					err := fr_index.Batch(batch)
					if err != nil {
						panic(err)
					}
					doneChan <- true
					progBar.FinishPrint("Finished Indexing!")
					return
				}
			}
		}
	}(fr_index)
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

func getNumFiles(dir string) int {
	count := 0
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count += 1
		}
		return err
	})
	return count
}
