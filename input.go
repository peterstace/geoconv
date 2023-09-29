package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func decodeInput(input []byte, format string, validate bool) (geom.Geometry, error) {
	switch format {
	case "any":
		return decodeUsingAny(input, validate)
	case "wkt":
		return decodeUsingWKT(input, validate)
	case "geojson":
		return decodeUsingGeoJSON(input, validate)
	case "lonlat":
		return decodeUsingLonLat(input, validate)
	case "tile":
		return decodeUsingTile(input, validate)
	case "seq":
		return decodeUsingSeq(input, validate)
	default:
		return geom.Geometry{}, fmt.Errorf("unknown geometry format: %v", format)
	}
}

func decodeUsingAny(input []byte, validate bool) (geom.Geometry, error) {
	for _, fn := range []func([]byte, bool) (geom.Geometry, error){
		decodeUsingWKT,
		decodeUsingGeoJSON,
		decodeUsingLonLat,
		decodeUsingTile,
		decodeUsingSeq,
	} {
		g, err := fn(input, validate)
		if err != nil {
			continue // try the next format
		}
		return g, nil
	}
	return geom.Geometry{}, errors.New("could not parse using any format")
}

func decodeUsingWKT(input []byte, validate bool) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	var nv []geom.NoValidate
	if !validate {
		nv = append(nv, geom.NoValidate{})
	}
	return geom.UnmarshalWKT(string(input), nv...)
}

func decodeUsingGeoJSON(input []byte, validate bool) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	var nv []geom.NoValidate
	if !validate {
		nv = append(nv, geom.NoValidate{})
	}
	return geom.UnmarshalGeoJSON(input, nv...)
}

func decodeUsingLonLat(input []byte, validate bool) (geom.Geometry, error) {
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
	pt := geom.XY{X: x, Y: y}.AsPoint()
	err = pt.Validate()
	return pt.AsGeometry(), err
}

func decodeUsingTile(input []byte, validate bool) (geom.Geometry, error) {
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

	env, err := Tile{z, x, y}.AsEnvelope()
	return env.AsGeometry(), err
}

func decodeUsingSeq(input []byte, validate bool) (geom.Geometry, error) {
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
		pt := seq.GetXY(0).AsPoint()
		return pt.AsGeometry(), pt.Validate()
	default:
		ring := geom.NewLineString(seq)
		poly := geom.NewPolygon([]geom.LineString{ring})
		return poly.AsGeometry(), poly.Validate()
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
