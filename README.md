# geoconv

Converts common geospatial data formats from one type to another. 
For example geojson into KML OR WKT into geojson etc.

## Supported Formats

Inputs to geoconv are taken from stdin.

If the `--input` flag is set, then its value is used as the input type code. 
If the flag is missing, then autodetection is used.

The following input types are supported:

| Type             | Code      | Example                                                   | Alternate separators |
| ---              | ---       | ---                                                       | ---                  |
| WKT              | `wkt`     | `POINT(151.2 -33.9)`                                      |                      |
| GeoJSON          | `geojson` | `{"type":"Point","coordinates":[151.2,-33.9]}`            |                      |
| Sequence         | `seq`     | `151.02,-33.45,150.61,-34.16,151.76,-34.15,151.02,-33.45` | comma, space         |
| Tile Coordinates | `tile`    | `21 1929379 1258703`                                      | space, forward slash |

If the `--output` flag is set, then its value is used as the output type code.
If the flag is missing, then all outputs are shown.

## Open geojson.io

If the `--show` flag is provided, then [geojson.io](geojson.io) will be opened
with the geometry from the input loaded.



## Installation

Either compile from source or install only the binary via
```
go install github.com/peterstace/geoconv@latest
``` 


## Usage

Pipe a valid geospatial data to the binary

1 ) Autodetect input format and convert to all supported formats

```
echo '{"type":"Point","coordinates":[151.2,-33.9]}' | geoconv
```

2 ) Explicitly set both input and output formats
```
echo 'POINT(151.2 -33.9)' | geoconv --input wkt --output geojson
```
