package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	MetaJsonPath string = IndexDir + "fdb_meta.json"
	IndexLimit   int    = 10
)

type IndexRecord map[string]string //`json:"index_record"`

func NewIndexRecord() IndexRecord {
	return IndexRecord(make(map[string]string))
}

func (ir IndexRecord) AddIndex(idx string) {
	ir[idx] = time.Now().Format(time.RFC3339)
}

func GetIndexRecord() (IndexRecord, error) {
	if _, err := os.Stat(MetaJsonPath); os.IsNotExist(err) {
		errDir := os.Mkdir(IndexDir, 0700)
		_, errFile := os.Create(MetaJsonPath)
		idxRecord := NewIndexRecord()
		if errDir != nil || errFile != nil {
			return nil, err
		}
		return idxRecord, nil
	}
	buffer, err := ioutil.ReadFile(MetaJsonPath)
	if err != nil {
		return nil, err
	}
	record := make(map[string]string)
	err = json.Unmarshal(buffer, &record)
	if err != nil {
		return nil, err
	}
	return IndexRecord(record), nil
}

func (ir IndexRecord) SaveToJson() error {
	res, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(MetaJsonPath, res, os.ModeType)
	if err != nil {
		panic(err)
		return err
	}
	return nil
}

func (ir IndexRecord) DirHasValidIndex(path string) bool {
	dirList := strings.Split(path, "/")
	// check if path is ever indexed (by itself or parent dirs)
	var indexPath string = ""
	for i := 0; i <= len(dirList); i++ {
		temp := getIndexName(strings.Join(dirList[:i], "/"))
		if _, ok := ir[temp]; ok {
			indexPath = temp
			break
		}
	}
	if indexPath == "" {
		return false
	}
	// check if index is young enough (indexed less than 6 hours ago)
	indexTime, _ := time.Parse(time.RFC3339, ir[indexPath])
	interval := time.Now().Sub(indexTime)
	if interval >= 6*60*time.Minute {
		return false
	}
	return true
}
