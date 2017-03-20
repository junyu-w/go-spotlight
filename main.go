package main

import (
	"github.com/DrakeW/go-spotlight/client"
	"os"
)

func main() {
	cli := client.GetCliApp()
	cli.Run(os.Args)

	//count := 1000
	//bar := pb.New(count)
	//bar.Start()
	//for i := 0; i < count; i++ {
	//	bar.Increment()
	//	time.Sleep(time.Millisecond)
	//}
	//bar.FinishPrint("The End!")

}
