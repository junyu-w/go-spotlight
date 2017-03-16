package models

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"strings"
)

const IndexDir string = "/Users/junyuwang/Desktop/fdb_idx/"

func newFrIndex(idxName string) bleve.Index {
	idxMapping := frIndexMapping()
	fr_idx, err := bleve.New(getIndexPath(idxName), idxMapping)
	if err != nil {
		panic(err)
	}
	return fr_idx
}

func frIndexMapping() *mapping.IndexMappingImpl {
	frIndexMapping := bleve.NewIndexMapping()

	fileMapping := bleve.NewDocumentMapping()
	frIndexMapping.AddDocumentMapping("file_record", fileMapping)

	recentContentMapping := bleve.NewTextFieldMapping()
	recentContentMapping.Analyzer = "en"
	fileMapping.AddFieldMappingsAt("RecentContent", recentContentMapping)

	fileNameMapping := bleve.NewTextFieldMapping()
	recentContentMapping.Analyzer = "standard"
	fileMapping.AddFieldMappingsAt("Name", fileNameMapping)

	timeMapping := bleve.NewDateTimeFieldMapping()
	fileMapping.AddFieldMappingsAt("AccessTime", timeMapping)
	fileMapping.AddFieldMappingsAt("ModifyTime", timeMapping)
	fileMapping.AddFieldMappingsAt("ChangeTime", timeMapping)

	frIndexMapping.TypeField = "type"
	frIndexMapping.DefaultAnalyzer = "en"

	return frIndexMapping
}

func getIndexName(dir string) string {
	idxName := strings.Join(strings.Split(dir, "/"), "_") + ".bleve"
	return idxName
}

func getIndexPath(idxName string) string {
	return IndexDir + idxName
}

func GetFrIndex(cwd string) (bleve.Index, error) {
	var fr_index bleve.Index
	// check idx record
	idxRecord, err := GetIndexRecord()
	if err != nil {
		panic(err)
		return nil, err
	}
	if index, valid := idxRecord.DirHasValidIndex(cwd); valid == false {
		fmt.Println("creating new index...")
		if index != "" {
			idxRecord.RemoveIndex(cwd)
		}
		idxName := getIndexName(cwd)
		fr_index = newFrIndex(idxName)
		idxRecord.AddIndex(idxName)
		idxRecord.SaveToJson()
	} else {
		fmt.Println("querying existed index...")
		fr_index, err = bleve.Open(IndexDir + index)
		if err != nil {
			panic(err)
		}
	}
	return fr_index, nil
}
