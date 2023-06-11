PYTHON ?= hatch run python

run-examples: 00 01 02

# omit top level parameters
00:
	$(PYTHON) simplify/__init__.py _examples/00parameters/main.yaml --format yaml > _examples/00parameters/simplified.yaml
.PHONY: 00
# omit $ref (exclude components)
01:
	$(PYTHON) simplify/__init__.py _examples/01deref-parameter/main.yaml --format yaml > _examples/01deref-parameter/simplified.yaml
.PHONY: 01
# omit allOf
02:
	$(PYTHON) simplify/__init__.py _examples/02allof/main.yaml --format yaml > _examples/02allof/simplified.yaml
.PHONY: 02


# OAPI_VERSION ?= v0.118.0
OAPI_VERSION ?= latest
validate-examples:
	go run github.com/getkin/kin-openapi/cmd/validate@$(OAPI_VERSION) ./_examples/00*/main.yaml
	go run github.com/getkin/kin-openapi/cmd/validate@$(OAPI_VERSION) ./_examples/01*/main.yaml
	go run github.com/getkin/kin-openapi/cmd/validate@$(OAPI_VERSION) ./_examples/02*/main.yaml
.PHONY: validate-examples