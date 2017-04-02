package client

import (
	"bytes"
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/gosuri/uitable"
	"github.com/logrusorgru/aurora"
	"io"
	"log"
	"os"
	"os/exec"
)

func formatSearchResult(result *bleve.SearchResult, showField bool) string {
	table := uitable.New()
	table.MaxColWidth = 120
	table.Wrap = true
	summary := aurora.Brown(fmt.Sprintf("%d matches, showing %d through %d, took %s\n", result.Total, result.Request.From+1, result.Request.From+len(result.Hits), result.Took))
	table.AddRow("")
	for _, hit := range result.Hits {
		path := hit.ID
		table.AddRow(aurora.Green("Path:"), path)
		if showField {
			for field, v := range hit.Fragments {
				if field != "Path" && field != "Extension" && field != "Name" {
					table.AddRow(fmt.Sprintf("%s:", aurora.Green(field)), fmt.Sprintf("...%s...", v))
				}
			}
			table.AddRow("") // blank
		}
	}
	table.AddRow("") // blank
	return summary.String() + table.String()
}

func pipeResultToPeco(result string) {
	cmd := exec.Command("peco", "--initial-index", "2")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Panic(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Panic(err)
	}
	err = cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	defer stdin.Close()
	io.Copy(stdin, bytes.NewBufferString(result))
	io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}
}
