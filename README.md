# geoconv

## Input

Inputs to geoconv are taken from stdin.

If the `-i` flag is set, then its value is used at the input type. If the flag
is missing, then autodetection is used.

The following input types are supported:

| Type             | Code      | Example                                        | Alternate separators |
| ---              | ---       | ---                                            | ---                  |
| WKT              | `wkt`     | `POINT(151.2 -33.9)`                           |                      |
| GeoJSON          | `geojson` | `{"type":"Point","coordinates":[151.2,-33.9]}` |                      |
| Lon Lat          | `lonlat`  | `151.2,33.9`                                   | comma, space         |
| Tile Coordinates | `tile`    | `21 1929379 1258703`                           | space, forward slash |
