# geoconv

## Formats

Inputs to geoconv are taken from stdin.

If the `-i` flag is set, then its value is used as the input type code. If the flag
is missing, then autodetection is used.

The following input types are supported:

| Type             | Code      | Example                                                   | Alternate separators |
| ---              | ---       | ---                                                       | ---                  |
| WKT              | `wkt`     | `POINT(151.2 -33.9)`                                      |                      |
| GeoJSON          | `geojson` | `{"type":"Point","coordinates":[151.2,-33.9]}`            |                      |
| Sequence         | `seq`     | `151.02,-33.45,150.61,-34.16,151.76,-34.15,151.02,-33.45` | comma, space         |
| Tile Coordinates | `tile`    | `21 1929379 1258703`                                      | space, forward slash |

If the `-t` flag is set, then its value is sued as the output type code. If the
flag is missing, then all outputs are shown.

# Open geojson.io

If the `-o` flag is provided, then [geojson.io](geojson.io) will be opened with
the geometry from the input loaded.

              
