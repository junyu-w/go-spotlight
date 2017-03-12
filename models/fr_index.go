package models

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"os"
)

const IndexName string = "file_db.bleve"

func newFrIndex() bleve.Index {
	idxMapping := frIndexMapping()
	fr_idx, err := bleve.New(IndexName, idxMapping)
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

func GetFrIndex() bleve.Index {
	var fr_index bleve.Index
	if _, err := os.Stat(IndexName); os.IsNotExist(err) {
		fmt.Println("creating new index...")
		fr_index = newFrIndex()

	} else {
		fmt.Println("querying existed index...")
		fr_index, err = bleve.Open(IndexName)
		if err != nil {
			panic(err)
		}
	}
	return fr_index
}
