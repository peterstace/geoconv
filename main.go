package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	inputFormat := flag.String("input", "any", "input format")
	outputFormat := flag.String("output", "all", "output format")
	disableValidation := flag.Bool("disable-validation", false, "disable validation")
	open := flag.Bool("show", false, "show in browser (geojson.io)")
	flag.Parse()

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("could not read stdin: %v", err)
	}

	inputGeom, err := decodeInput(input, *inputFormat, *disableValidation)
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
