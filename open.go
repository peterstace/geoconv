package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"os/exec"
	"strings"

	"github.com/peterstace/simplefeatures/geom"
)

func openInBrowser(g geom.Geometry) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(g); err != nil {
		return err
	}

	u := url.URL{
		Scheme:   "http",
		Host:     "geojson.io",
		Path:     "/",
		Fragment: "data=data:application/json," + strings.TrimSpace(buf.String()),
	}

	out, err := exec.Command("open", u.String()).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
