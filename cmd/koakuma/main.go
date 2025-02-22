package main

import (
	"flag"
	"gokoakumachan/koakumachan"
)

func main() {
	confFilename := flag.String("conf", "config.yaml", "config file")
	flag.Parse()

	err := koakumachan.Main(*confFilename)
	if err != nil {
		panic(err)
	}
}
