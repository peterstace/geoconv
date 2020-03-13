package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	format := flag.String("i", "any", "input format")
	flag.Parse()

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read stdin: %v", err)
	}

	inputGeom, err := decodeInput(input, *format)
	if err != nil {
		log.Fatalf("decoding input: %v", err)
	}

	output(inputGeom)
}
