package main

import (
	"math"

	"github.com/peterstace/simplefeatures/geom"
)

type Tile struct {
	Z, X, Y int
}

func LonLatToTile(lon, lat float64, z int) Tile {
	latrad := lat * math.Pi / 180
	tanlat := math.Tan(latrad)
	coslat := math.Cos(latrad)
	twoPowZ := math.Exp2(float64(z))
	x := int((lon + 180.0) / 360.0 * twoPowZ)
	y := int((1.0 - math.Log(tanlat+1.0/coslat)/math.Pi) * 0.5 * twoPowZ)
	return Tile{Z: z, X: x, Y: y}
}

func (t Tile) AsEnvelope() geom.Envelope {
	return geom.NewEnvelope(
		geom.XY{
			X: tileCoordsToLon(t.Z, t.X),
			Y: tileCoordsToLat(t.Z, t.Y),
		},
		geom.XY{
			X: tileCoordsToLon(t.Z, t.X+1),
			Y: tileCoordsToLat(t.Z, t.Y+1),
		},
	)
}

func tileCoordsToLon(z, x int) float64 {
	return float64(x)/float64(int(1)<<z)*360 - 180
}

func tileCoordsToLat(z, y int) float64 {
	return 180 / math.Pi * math.Atan(math.Sinh(math.Pi*(1-2*(float64(y)/float64(int(1)<<z)))))
}
