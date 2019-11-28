package main

import (
	"encoding/json"
	"math"
	"os"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func main() {
	c := TileCoordinates{21, 478841, 863802}
	g := c.ToEnvelope().AsGeometry()
	json.NewEncoder(os.Stdout).Encode(g)

	g2, err := geom.UnmarshalWKT(strings.NewReader("POLYGON((115.904195 -31.931195,115.904195 -31.909139,115.87394 -31.909139,115.87394 -31.931195,115.904195 -31.931195))"))
	if err != nil {
		panic(err)
	}
	json.NewEncoder(os.Stdout).Encode(g2)
}

type TileCoordinates struct {
	Z, X, Y int
}

func (c TileCoordinates) ToEnvelope() geom.Envelope {
	n := 1 << c.Z
	lonA := float64(c.X)/float64(n)*360 - 180
	latA := 180 / math.Pi * math.Atan(math.Sinh(math.Pi*(1-2*float64(c.Y)/float64(n))))
	lonB := float64(c.X+1)/float64(n)*360 - 180
	latB := 180 / math.Pi * math.Atan(math.Sinh(math.Pi*(1-2*float64(c.Y+1)/float64(n))))
	return geom.NewEnvelope(
		geom.XY{X: lonA, Y: latA},
		geom.XY{X: lonB, Y: latB},
	)
}

// PolySequence is a comma separated list of scalars that represents a polygon
// without any holes. The first, third, fifth etc elements are longitudes, and
// the second, fourth etc elements are latitudes.
//
// E.g.
//
// 115.904195,-31.931195,115.904195,-31.909139,115.87394,-31.909139,115.87394,-31.931195,115.904195,-31.931195
type PolySequence struct {
	Coordinates []float64
}

func (s PolySequence) ToGeometry() geom.Envelope {
}

/*

n = 2 ^ zoom
lon_deg = xtile / n * 360.0 - 180.0
lat_rad = arctan(sinh(π * (1 - 2 * ytile / n)))
lat_deg = lat_rad * 180.0 / π

*/
