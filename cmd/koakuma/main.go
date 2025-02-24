package main

import (
	"flag"
	"log"

	"github.com/mudream4869/gokoakumachan/koakumachan"
)

func main() {
	confFilename := flag.String("conf", "config.yaml", "config file")
	flag.Parse()

	err := koakumachan.Main(*confFilename)
	if err != nil {
		log.Fatalln(err)
	}
}
