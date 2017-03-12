package main

import (
	//"fmt"
	"github.com/DrakeW/FileDB/client"
	//"github.com/blevesearch/bleve"
	"os"
	//"time"
)

func main() {
	cli := client.GetCliApp()
	cli.Run(os.Args)
	//// time query
	//// TODO: time isn't working correctly
	//begin, _ := time.Parse(time.RFC3339, "2018-03-01T00:00:00Z-7:00")
	//end, _ := time.Parse(time.RFC3339, "2018-03-11T24:00:05Z-7:00")
	//query := bleve.NewDateRangeQuery(begin, end)
	//searchRequest := bleve.NewSearchRequest(query)
	//searchResult, _ := fr_index.Search(searchRequest)
	//fmt.Println(searchResult)
	//
	//// content query
	//query2 := bleve.NewQueryStringQuery("test file")
	//searchRequest2 := bleve.NewSearchRequest(query2)
	//searchResult2, err := fr_index.Search(searchRequest2)
	//// TODO: ansi highlighter is not registered
	//fmt.Println(err)
	//fmt.Println(searchResult2)
}
