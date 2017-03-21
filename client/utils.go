package client

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gosuri/uitable"
)

func formatSearchResult(result *bleve.SearchResult) string {
	table := uitable.New()
	table.MaxColWidth = 80
	table.Wrap = true
	summary := fmt.Sprintf("%d matches, showing %d through %d, took %s\n", result.Total, result.Request.From+1, result.Request.From+len(result.Hits), result.Took)
	table.AddRow("")
	for _, hit := range result.Hits {
		path := hit.ID
		table.AddRow("Path:", path)
		for field, v := range hit.Fragments {
			table.AddRow(fmt.Sprintf("%s:", field), v)
		}
		table.AddRow("") // blank
	}
	return summary + table.String()
}
