package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/peterstace/simplefeatures/geom"
)

func main() {
	inputFormat := flag.String("input", "any", "input format")
	outputFormat := flag.String("output", "all", "output format")
	disableValidation := flag.Bool("disable-validation", false, "disable validation")
	outputPrecisionStr := flag.String("dp", "", "output precision, an integer number of decimal places")
	boundingBox := flag.Bool("bbox", false, "output bounding box instead of original geometry")
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
	if *boundingBox {
		inputGeom = inputGeom.Envelope().AsGeometry()
	}

	if *outputPrecisionStr != "" {
		dp, err := strconv.Atoi(*outputPrecisionStr)
		if err != nil {
			log.Fatalf("could not parse output precision as int: %v", err)
		}
		factor := math.Pow(10, float64(dp))
		inputGeom, err = inputGeom.TransformXY(func(xy geom.XY) geom.XY {
			xy.X = math.Round(xy.X*factor) / factor
			xy.Y = math.Round(xy.Y*factor) / factor
			return xy
		})
		if err != nil {
			log.Fatalf("could not modify precision: %v", err)
		}
	}

	output(inputGeom, *outputFormat)

	if *open {
		if err := openInBrowser(inputGeom); err != nil {
			log.Fatalf("could not open in browser: %v", err)
		}
	}
}
