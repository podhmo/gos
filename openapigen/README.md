# openapigen

generate the subset of openapi by builder

## examples

see: [examples](https://github.com/podhmo/gos/tree/main/openapigen/_examples)

## getting started

```console
$ mkdir -p tools
$ go run github.com/podhmo/gos/openapigen/cmd/init@latest > tools/gen.go
$ go run tools/gen.go > openapi.json
```