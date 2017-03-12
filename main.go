package main

import (
	"fmt"
	"github.com/DrakeW/FileDB/models"
	"github.com/blevesearch/bleve"
	"os"
	"time"
)

func main() {
	var fr_index bleve.Index
	if _, err := os.Stat(models.IndexName); os.IsNotExist(err) {
		// path/to/whatever does not exist
		fmt.Println("creating new index...")
		fr_index = models.NewFrIndex()

	} else {
		fmt.Println("opening index...")
		fr_index, err = bleve.Open(models.IndexName)
		if err != nil {
			panic(err)
		}
	}
	defer fr_index.Close()
	//models.IndexAllFiles("../../../", &fr_index)
	// time query
	// TODO: time isn't working correctly
	begin, _ := time.Parse(time.RFC3339, "2018-03-01T00:00:00Z-7:00")
	end, _ := time.Parse(time.RFC3339, "2018-03-11T24:00:05Z-7:00")
	query := bleve.NewDateRangeQuery(begin, end)
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, _ := fr_index.Search(searchRequest)
	fmt.Println(searchResult)

	// content query
	query2 := bleve.NewQueryStringQuery("test file")
	searchRequest2 := bleve.NewSearchRequest(query2)
	//searchRequest2.Highlight = bleve.NewHighlightWithStyle("ansi")
	//searchRequest2.Highlight.AddField("RecentContent")
	searchResult2, _ := fr_index.Search(searchRequest2)
	// TODO: ansi highlighter is not registered
	fmt.Println(searchResult2)
}
