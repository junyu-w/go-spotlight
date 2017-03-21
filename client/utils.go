package client

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gosuri/uitable"
	"github.com/logrusorgru/aurora"
)

func formatSearchResult(result *bleve.SearchResult) string {
	table := uitable.New()
	table.MaxColWidth = 120
	table.Wrap = true
	summary := aurora.Brown(fmt.Sprintf("%d matches, showing %d through %d, took %s\n", result.Total, result.Request.From+1, result.Request.From+len(result.Hits), result.Took))
	table.AddRow("")
	for _, hit := range result.Hits {
		path := hit.ID
		table.AddRow(aurora.Green("Path:"), path)
		for field, v := range hit.Fragments {
			if field != "Path" && field != "Extension" {
				table.AddRow(fmt.Sprintf("%s:", aurora.Green(field)), fmt.Sprintf("...%s...", v))
			}
		}
		table.AddRow("") // blank
	}
	return summary.String() + table.String()
}
