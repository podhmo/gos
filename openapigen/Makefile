run-example: 00 01 02 03 04
.PHONY: run-example

# simple schema only example
00:
	go run ./_examples/00schema > ./_examples/00schema/openapi.json
.PHONY: 00
# simple router example
01:
	go run ./_examples/01router > ./_examples/01router/openapi.json
.PHONY: 01
# separated files example
02:
	go run ./_examples/02separated > ./_examples/02separated/openapi.json
.PHONY: 02
# generate also go code
03:
	mkdir -p ./_examples/03tocode/testdata
	go run ./_examples/03tocode
.PHONY: 03
# generate pagination example (extend schema by a helper function)
04:
	go run ./_examples/04pagination > ./_examples/04pagination/openapi.json

_validate:
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/01router/openapi.json
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/02separated/openapi.json
.PHONY: _validate