package models

import (
	"encoding/json"
	"fmt"
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
		f, errFile := os.Create(MetaJsonPath)
		defer f.Close()
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

// check if dir or any parent dir was indexed
func (ir IndexRecord) GetUsableIndex(path string) string {
	dirList := strings.Split(path, "/")
	var indexName string = ""
	for i := 0; i <= len(dirList); i++ {
		temp := getIndexName(strings.Join(dirList[:i], "/"))
		if _, ok := ir[temp]; ok {
			indexName = temp
			break
		}
	}
	return indexName
}

func (ir IndexRecord) DirHasValidIndex(path string) (string, bool) {
	// check if path is ever indexed (by itself or parent dirs)
	indexName := ir.GetUsableIndex(path)
	if indexName == "" {
		return indexName, false
	}
	// check if index is young enough (indexed less than 6 hours ago)
	indexTime, _ := time.Parse(time.RFC3339, ir[indexName])
	interval := time.Now().Sub(indexTime)
	if interval >= 6*60*time.Minute {
		return indexName, false
	}
	return indexName, true
}

func (ir IndexRecord) RemoveIndex(indexName string) error {
	fmt.Println("cleaning old index...")
	err := os.RemoveAll(IndexDir + indexName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
