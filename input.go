package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func decodeInput(input []byte, format string) (geom.Geometry, error) {
	switch format {
	case "any":
		return decodeUsingAny(input)
	case "wkt":
		return decodeUsingWKT(input)
	case "geojson":
		return decodeUsingGeoJSON(input)
	case "lonlat":
		return decodeUsingLonLat(input)
	case "tile":
		return decodeUsingTile(input)
	default:
		return geom.Geometry{}, fmt.Errorf("unknown geometry format: %v", format)
	}
}

func decodeUsingAny(input []byte) (geom.Geometry, error) {
	for _, fn := range []func([]byte) (geom.Geometry, error){
		decodeUsingWKT,
		decodeUsingGeoJSON,
		decodeUsingLonLat,
		decodeUsingTile,
	} {
		g, err := fn(input)
		if err != nil {
			continue // try the next format
		}
		return g, nil
	}
	return geom.Geometry{}, errors.New("could not parse using any format")
}

func decodeUsingWKT(input []byte) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	return geom.UnmarshalWKT(bytes.NewReader(input))
}

func decodeUsingGeoJSON(input []byte) (geom.Geometry, error) {
	// TODO: If we could detect _geometry validation_ errors here, we could
	// bubble them up rather than using wrongFormatErr.
	var g geom.Geometry
	err := json.NewDecoder(bytes.NewReader(input)).Decode(&g)
	return g, err
}

func decodeUsingLonLat(input []byte) (geom.Geometry, error) {
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
	return geom.NewPointF(x, y).AsGeometry(), nil
}

func decodeUsingTile(input []byte) (geom.Geometry, error) {
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

func coalesce(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
