package main

import (
	"github.com/DrakeW/FileDB/client"
	"os"
)

func main() {
	cli := client.GetCliApp()
	cli.Run(os.Args)
}
