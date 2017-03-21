package client

import (
	"fmt"
	"github.com/DrakeW/go-spotlight/models"
	"github.com/blevesearch/bleve"
	_ "github.com/blevesearch/bleve/config"
	"github.com/blevesearch/bleve/search/query"
	"github.com/urfave/cli"
	"os"
)

// sample search: "fdb -k name -t -2d~0d -ext txt --hint "hello world good desk"

func GetCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = "File DB"
	app.Usage = "Smart file search engine"

	app.Action = func(c *cli.Context) error {
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:    "strictquery",
			Aliases: []string{"sq"},
			Usage:   "Search for files according to query input",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "time, t",
					Value: "",
					Usage: "date range for your search, eg. -3~0 means last 3 days (mandatory field)",
				}, cli.StringFlag{
					Name:  "extension, ext",
					Value: "",
					Usage: "file extension to limit search ranges (default empty)",
				}, cli.StringFlag{
					Name:  "words, w",
					Value: "",
					Usage: "Enter words you remember that you typed in this file (default empty)",
				},
			},
			Action: executeStrictQuery,
		}, {
			Name:    "fuzzyquery",
			Aliases: []string{"q"},
			Usage:   "just dump whatever in your mind to the query",
			Action:  executeFuzzyQuery,
		},
	}
	return app
}

func executeStrictQuery(c *cli.Context) error {
	timeRange := c.String("time")
	fileExtension := c.String("extension")
	hint := c.String("words")
	if timeRange == "" {
		fmt.Println("Please check help with \"fdb help sq\"")
		return nil
	}

	curDir, _ := os.Getwd()
	fr_index, err := models.GetFrIndex(curDir)
	if err != nil {
		return err
	}
	defer fr_index.Close()

	query := compileQuery(timeRange, fileExtension, hint)

	searchRequest := bleve.NewSearchRequestOptions(query, 10, 0, false)
	searchRequest.Highlight = bleve.NewHighlightWithStyle("ansi")
	searchResult, err := fr_index.Search(searchRequest)
	if err != nil {
		return err
	}
	fmt.Println(formatSearchResult(searchResult))
	return nil
}

func executeFuzzyQuery(c *cli.Context) error {
	curDir, _ := os.Getwd()
	fr_index, err := models.GetFrIndex(curDir)
	if err != nil {
		return err
	}
	defer fr_index.Close()

	queries := make([]query.Query, c.NArg(), c.NArg())
	for i := 0; i < c.NArg(); i++ {
		queries[i] = bleve.NewQueryStringQuery(c.Args()[i])
	}
	conjQuery := bleve.NewConjunctionQuery(queries...)

	searchRequest := bleve.NewSearchRequestOptions(conjQuery, 10, 0, false)
	searchRequest.Highlight = bleve.NewHighlightWithStyle("ansi")
	searchResult, err := fr_index.Search(searchRequest)
	if err != nil {
		return err
	}
	fmt.Println(formatSearchResult(searchResult))
	return nil
}
