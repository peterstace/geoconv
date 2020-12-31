package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	inputFormat := flag.String("i", "any", "input format")
	outputFormat := flag.String("t", "all", "output format")
	open := flag.Bool("o", false, "open in browser (geojson.io)")
	flag.Parse()

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read stdin: %v", err)
	}

	inputGeom, err := decodeInput(input, *inputFormat)
	if err != nil {
		log.Fatalf("decoding input: %v", err)
	}

	output(inputGeom, *outputFormat)

	if *open {
		if err := openInBrowser(inputGeom); err != nil {
			log.Fatalf("could not open in browser: %v", err)
		}
	}
}
