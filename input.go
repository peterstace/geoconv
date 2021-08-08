package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func decodeInput(input []byte, format string, disableValidation bool) (geom.Geometry, error) {
	var opts []geom.ConstructorOption
	if disableValidation {
		opts = []geom.ConstructorOption{geom.DisableAllValidations}
	}
	switch format {
	case "any":
		return decodeUsingAny(input, opts)
	case "wkt":
		return decodeUsingWKT(input, opts)
	case "geojson":
		return decodeUsingGeoJSON(input, opts)
	case "lonlat":
		return decodeUsingLonLat(input, opts)
	case "tile":
		return decodeUsingTile(input, opts)
	case "seq":
		return decodeUsingSeq(input, opts)
	default:
		return geom.Geometry{}, fmt.Errorf("unknown geometry format: %v", format)
	}
}

func decodeUsingAny(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	for _, fn := range []func([]byte, []geom.ConstructorOption) (geom.Geometry, error){
		decodeUsingWKT,
		decodeUsingGeoJSON,
		decodeUsingLonLat,
		decodeUsingTile,
		decodeUsingSeq,
	} {
		g, err := fn(input, opts)
		if err != nil {
			continue // try the next format
		}
		return g, nil
	}
	return geom.Geometry{}, errors.New("could not parse using any format")
}

func decodeUsingWKT(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	return geom.UnmarshalWKT(string(input), opts...)
}

func decodeUsingGeoJSON(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	return geom.UnmarshalGeoJSON(input, opts...)
}

func decodeUsingLonLat(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	parts := strings.Split(string(input), " ")
	if len(parts) != 2 {
		parts = strings.Split(string(input), ",")
		if len(parts) != 2 {
			return geom.Geometry{}, errors.New("could not extract 2 parts")
		}
	}
	parts[0] = strings.TrimSpace(parts[0])
	parts[1] = strings.TrimSpace(parts[1])

	x, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return geom.Geometry{}, err
	}
	y, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return geom.Geometry{}, err
	}
	return geom.NewPointFromXY(geom.XY{X: x, Y: y}).AsGeometry(), nil
}

func decodeUsingTile(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	parts := strings.Split(string(input), " ")
	if len(parts) != 3 {
		parts = strings.Split(string(input), "/")
		if len(parts) != 3 {
			return geom.Geometry{}, errors.New("could not extract 3 parts")
		}
	}
	for i := 0; i < 3; i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}

	z, errz := strconv.Atoi(parts[0])
	x, errx := strconv.Atoi(parts[1])
	y, erry := strconv.Atoi(parts[2])
	if err := coalesce(errz, errx, erry); err != nil {
		return geom.Geometry{}, err
	}

	return Tile{z, x, y}.AsEnvelope().AsGeometry(), nil
}

func decodeUsingSeq(input []byte, opts []geom.ConstructorOption) (geom.Geometry, error) {
	parts := strings.Split(string(input), ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	floats := make([]float64, len(parts))
	for i, str := range parts {
		var err error
		floats[i], err = strconv.ParseFloat(str, 64)
		if err != nil {
			return geom.Geometry{}, err
		}
	}

	if len(floats)%2 != 0 {
		return geom.Geometry{}, fmt.Errorf("number of items in sequence must be even")
	}
	seq := geom.NewSequence(floats, geom.DimXY)
	switch seq.Length() {
	case 0:
		return geom.Geometry{}, fmt.Errorf("no items in sequence")
	case 1:
		return geom.NewPointFromXY(seq.GetXY(0), opts...).AsGeometry(), nil
	default:
		ring, err := geom.NewLineString(seq, opts...)
		if err != nil {
			return geom.Geometry{}, err
		}
		poly, err := geom.NewPolygonFromRings([]geom.LineString{ring}, opts...)
		return poly.AsGeometry(), err
	}
}

func coalesce(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
