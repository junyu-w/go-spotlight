package models

import (
	"encoding/json"
	"os"
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

// TODO: finish index validity checking

func (ir *IndexRecord) DirHasValidIndex(dir string) bool {
	return ir.dirAlreadyIndexed(dir) && ir.dirIndexIsRecent(dir)
}

func (ir *IndexRecord) dirAlreadyIndexed(dir string) bool {
	return false
}

func (ir *IndexRecord) dirIndexIsRecent(dir string) bool {
	return false
}
