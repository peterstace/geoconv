package main

import (
	"strconv"
	"testing"
)

func TestDecodeInput(t *testing.T) {
	for i, tt := range []struct {
		input, format string
		wantErr       bool
	}{
		{"POINT(1 2)", "any", false},
		{`{"type":"Point","coordinates":[1,2]}`, "any", false},
		{"1 2", "any", false},
		{"1,2", "any", false},
		{"21 1929379 1258703", "any", false},

		{"POINT(1 2)", "wkt", false},
		{`{"type":"Point","coordinates":[1,2]}`, "wkt", true},
		{"1 2", "wkt", true},
		{"1,2", "wkt", true},
		{"21 1929379 1258703", "wkt", true},

		{"POINT(1 2)", "geojson", true},
		{`{"type":"Point","coordinates":[1,2]}`, "geojson", false},
		{"1 2", "geojson", true},
		{"1,2", "geojson", true},
		{"21 1929379 1258703", "geojson", true},

		{"POINT(1 2)", "lonlat", true},
		{`{"type":"Point","coordinates":[1,2]}`, "lonlat", true},
		{"1 2", "lonlat", false},
		{"1,2", "lonlat", false},
		{"21 1929379 1258703", "lonlat", true},

		{"POINT(1 2)", "tile", true},
		{`{"type":"Point","coordinates":[1,2]}`, "tile", true},
		{"1 2", "tile", true},
		{"1,2", "tile", true},
		{"21 1929379 1258703", "tile", false},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			_, err := decodeInput([]byte(tt.input), tt.format, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr: %v, err: %v", tt.wantErr, err)
			}
		})
	}
}
