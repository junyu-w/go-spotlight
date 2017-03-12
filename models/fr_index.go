package models

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
)

const IndexName string = "file_db.bleve"

func NewFrIndex() bleve.Index {
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
