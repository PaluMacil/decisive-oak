package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(flag.Args()) < 1 {
		fmt.Println("No file given for analysis")
		os.Exit(1)
	}
}


