package main

import (
	"github.com/DrakeW/fileDB/client"
	"os"
)

func main() {
	cli := client.GetCliApp()
	cli.Run(os.Args)
}
