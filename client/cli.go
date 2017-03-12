package client

import (
	"fmt"
	"github.com/DrakeW/FileDB/models"
	"github.com/blevesearch/bleve"
	"github.com/urfave/cli"
	"os"
)

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
			Action: func(c *cli.Context) error {
				return nil
			},
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
	var fr_index bleve.Index
	if _, err := os.Stat(models.IndexName); os.IsNotExist(err) {
		fmt.Println("creating new index...")
		fr_index = models.NewFrIndex()

	} else {
		fmt.Println("opening index...")
		fr_index, err = bleve.Open(models.IndexName)
		if err != nil {
			return err
		}
	}
	defer fr_index.Close()
	models.IndexAllFiles(dir, &fr_index)
	return nil
}
