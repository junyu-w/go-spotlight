package client

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gosuri/uitable"
)

func formatSearchResult(result *bleve.SearchResult) *uitable.Table {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	for _, hit := range result.Hits {
		path := hit.ID
		table.AddRow("Path:", path)
		for field, v := range hit.Fragments {
			table.AddRow(fmt.Sprintf("%s:", field), v)
		}
	}
	return table
}
