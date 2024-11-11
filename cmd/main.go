package main

import (
	"log"

	"github.com/jborkows/gotemplate/internal/example"
	"github.com/jborkows/gotemplate/internal/logs"
)

func main() {
	logger, err := logs.Initialize(logs.FileLogger("gotemplate.log"))
	if err != nil {
		panic(err)
	}
	log.Println("Starting gotemplate")
	log.Println("1 + 2 = ", example.Example(1, 2))
	defer logger.Close()

}
