package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/peterstace/simplefeatures/geom"
)

func output(g geom.Geometry) {
	fmt.Println("WKT")
	fmt.Println(g.AsText())
	fmt.Println()

	fmt.Println("GeoJSON")
	json.NewEncoder(os.Stdout).Encode(g)
	fmt.Println()

	centroid, ok := g.Centroid().XY()
	if !ok {
		return
	}

	area := g.Area()
	if area > 0 {
		fmt.Println("Area (square meters)")
		earthDiameterMeters := 40075161.2
		metersPerDegree := earthDiameterMeters / 360 * math.Cos(centroid.Y*math.Pi/180)
		fmt.Printf("%.1f \n\n", area*metersPerDegree*metersPerDegree)
	}

	fmt.Println("Tiles")
	for z := 15; z <= 22; z++ {
		tile := LonLatToTile(centroid.X, centroid.Y, z)
		fmt.Printf("%d/%d/%d\n", tile.Z, tile.X, tile.Y)
	}
}
