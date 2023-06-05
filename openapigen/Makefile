run-example: 00schema 01router 02separated
.PHONY: run-example

# simple schema only example
00schema:
	go run ./_examples/00schema > ./_examples/00schema/openapi.json
.PHONY: 00schema
# simple router example
01router:
	go run ./_examples/01router > ./_examples/01router/openapi.json
.PHONY: 01router
# separated files example
02separated:
	go run ./_examples/02separated > ./_examples/02separated/openapi.json
.PHONY: 02separated
03tocode:
	mkdir -p ./_examples/03tocode/testdata
	go run ./_examples/03tocode > ./_examples/03tocode/testdata/generated.go
.PHONY: 03tocode

_validate:
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/01router/openapi.json
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/02separated/openapi.json
.PHONY: _validate