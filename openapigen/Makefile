run-example: 00schema 01router _default

_default:
	go run ./_examples/ > ./_examples/openapi.json

# simple schema only example
00schema:
	go run ./_examples/00schema > ./_examples/00schema/openapi.json
# simple router example
01router:
	go run ./_examples/01router > ./_examples/01router/openapi.json

.PHONY: run-example 00schema 01router _default