package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/haya14busa/goplay"
)

func main() {
	flag.Parse()

	code := os.Stdin
	if len(flag.Args()) > 0 {
		path := flag.Arg(0)
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()
		code = file
	}

	if err := goplay.DefaultClient.Run(code, os.Stdout, os.Stderr); err != nil {
		fmt.Println(err)
	}
}
