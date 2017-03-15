package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	JsonFilePath string = "./fdb_index_record.json"
	IndexLimit   int    = 10
)

type IndexRecord struct {
	indexRecord map[string]string    `json:"index_record"`
	timeRecord  map[string]time.Time `json:"time_record"`
	indexCount  int                  `json:"index_count"`
}

func GetIndexRecord() (*IndexRecord, error) {
	json_f, err := os.Open(JsonFilePath)
	if err != nil {
		return nil, err
	}
	defer json_f.Close()
	buffer := make([]byte, 0)
	_, err = json_f.Read(buffer)
	if err != nil {
		return nil, err
	}
	record := &IndexRecord{}
	err = json.Unmarshal(buffer, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (ir *IndexRecord) SaveToJson() error {
	res, err := json.Marshal(ir)
	if err != nil {
		return err
	}
	json_f, err := os.Open(JsonFilePath)
	if err != nil {
		return err
	}
	defer json_f.Close()
	_, err = json_f.Write(res)
	if err != nil {
		return err
	}
	return nil
}

func (ir *IndexRecord) DirHasValidIndex(path string) bool {
	dirList := strings.Split("/", path)
	// check if path is ever indexed (by itself or parent dirs)
	var indexPath string = ""
	for i := 0; i < len(dirList); i++ {
		temp := filepath.Join(dirList[:i]...)
		if _, ok := ir.indexRecord[temp]; ok {
			indexPath = temp
			break
		}
	}
	if indexPath == "" {
		return false
	}
	// check if index is young enough (indexed less than 6 hours ago)
	interval := time.Now().Sub(ir.timeRecord[indexPath])
	if interval >= 6*60*time.Minute {
		return false
	}
	return true
}
