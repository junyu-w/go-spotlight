package models

import (
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"os/user"
	"strings"
)

var IndexDir string = getHomeDir() + "/.fdb_idx/"

func getHomeDir() string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return user.HomeDir
}

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
		idxName := getIndexName(cwd)
		if index != "" {
			idxRecord.RemoveIndex(index)
		}
		fr_index = newFrIndex(idxName)
		// index files inside cwd
		idxRecord.AddIndex(idxName)
		idxRecord.SaveToJson()
		doneChan := make(chan bool)
		StartIndexing(cwd, fr_index, doneChan)
		<-doneChan
	} else {
		fr_index, err = bleve.Open(IndexDir + index)
		if err != nil {
			panic(err)
		}
	}
	return fr_index, nil
}
