package main

import (
	"github.com/DrakeW/go-spotlight/client"
	"os"
)

func main() {
	cli := client.GetCliApp()
	cli.Run(os.Args)
}
