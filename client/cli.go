package client

import (
	"fmt"
	"github.com/DrakeW/FileDB/models"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/urfave/cli"
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
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Index files under specified directory recrusively",
			Action:  executeIndexCommand,
		},
		{
			Name:    "query",
			Aliases: []string{"q"},
			Usage:   "Search for files according to query input",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "keyword, k",
					Value: "path",
					Usage: "searching for file 'name' or file 'path'",
				}, cli.StringFlag{
					Name:  "time, t",
					Value: "-1d~0d",
					Usage: "date range for your search",
				}, cli.StringFlag{
					Name:  "extension, ext",
					Value: "_all",
					Usage: "file extension to limit search ranges",
				}, cli.StringFlag{
					Name:  "hint",
					Value: "_None",
					Usage: "Enter any word you remember that you typed in this file, it can be inaccurate and fuzzy",
				},
			},
			Action: executeStrictQuery,
		}, {
			Name:    "fuzzyQuery",
			Aliases: []string{"fq"},
			Usage:   "just dump whatever in your mind to the query",
			Action:  executeFuzzyQuery,
		},
	}
	return app
}

func executeIndexCommand(c *cli.Context) error {
	args := c.Args()
	var dir string
	if len(args) > 0 {
		dir = args[0]
	} else {
		dir = "./"
	}
	fr_index := models.GetFrIndex()
	defer fr_index.Close()
	models.IndexAllFiles(dir, &fr_index)
	return nil
}

func executeStrictQuery(c *cli.Context) error {
	//keyword := c.String("keyword")
	timeRange := c.String("time")
	fileExtension := c.String("extension")
	hint := c.String("hint")

	query := compileQuery(timeRange, fileExtension, hint)
	fr_index := models.GetFrIndex()
	defer fr_index.Close()
	searchRequest := bleve.NewSearchRequest(query)
	searchResult, err := fr_index.Search(searchRequest)
	if err != nil {
		return err
	}
	// TODO: ansi highlighter is not registered
	fmt.Println(searchResult)
	return nil
}

func executeFuzzyQuery(c *cli.Context) error {
	fr_index := models.GetFrIndex()
	defer fr_index.Close()

	queries := make([]query.Query, c.NArg(), c.NArg())
	for i := 0; i < c.NArg(); i++ {
		queries[i] = bleve.NewQueryStringQuery(c.Args()[i])
	}
	conjQuery := bleve.NewConjunctionQuery(queries...)

	searchRequest := bleve.NewSearchRequest(conjQuery)
	searchResult, err := fr_index.Search(searchRequest)
	if err != nil {
		return err
	}
	fmt.Println(searchResult)
	return nil
}
