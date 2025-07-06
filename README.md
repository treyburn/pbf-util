# pbf-util
Command line utility for working with pbf encoded mvt tiles.

## Commands
### Decode
Decode allows for reading a pbf encoded mvt tile and pretty prints the json contents to stdout. Gzip deflating is handled by default, if necessary.

```shell
pmtiles tile path/to/file.pmtiles z x y | pbf-util decode 
```

Works nicely with `jq` as well:
```shell
pmtiles tile path/to/file.pmtiles z x y | pbf-util decode  | jq 'keys'
```

## Installation
Installing this command line utility requires `go` version 1.18 or newer.

Ensure that your `$GOBIN` or `$GOPATH/bin` are in your path.
```shell
export PATH="$(go env GOBIN):$(go env GOPATH)/bin:$PATH"
```

Once exported, you can simply `go install`.
```shell
go install github.com/treyburn/pbf-util
```