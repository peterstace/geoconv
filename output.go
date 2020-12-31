package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func output(g geom.Geometry, outputFormat string) {
	switch outputFormat {
	case "wkt":
		fmt.Println(g.AsText())
	case "geojson":
		json.NewEncoder(os.Stdout).Encode(g)
	case "seq":
		if g.IsPolygon() {
			poly := g.AsPolygon()
			seqs := []geom.Sequence{poly.ExteriorRing().Coordinates()}
			for i := 0; i < poly.NumInteriorRings(); i++ {
				seqs = append(seqs, poly.InteriorRingN(i).Coordinates())
			}

			var coords []string
			for _, seq := range seqs {
				for i := 0; i < seq.Length(); i++ {
					xy := seq.GetXY(i)
					coords = append(coords,
						strconv.FormatFloat(xy.X, 'f', -1, 64),
						strconv.FormatFloat(xy.Y, 'f', -1, 64),
					)
				}
			}
			fmt.Println(strings.Join(coords, ","))
		}
	case "all":
		fmt.Println("WKT")
		fmt.Println(g.AsText())
		fmt.Println()

		fmt.Println("GeoJSON")
		json.NewEncoder(os.Stdout).Encode(g)
		fmt.Println()

		if g.IsPolygon() && g.AsPolygon().NumInteriorRings() == 0 {
			seq := g.AsPolygon().ExteriorRing().Coordinates()
			var coords []string
			for i := 0; i < seq.Length(); i++ {
				xy := seq.GetXY(i)
				coords = append(coords,
					strconv.FormatFloat(xy.X, 'f', -1, 64),
					strconv.FormatFloat(xy.Y, 'f', -1, 64),
				)
			}
			fmt.Println("Sequence")
			fmt.Println(strings.Join(coords, ","))
			fmt.Println()
		}

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
}
