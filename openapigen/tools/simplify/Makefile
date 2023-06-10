PYTHON ?= hatch run python

run-examples: 00

00:
	$(PYTHON) simplify/__init__.py _examples/00parameters/main.yaml > _examples/00parameters/simplified.yaml
.PHONY: 00
01:
	$(PYTHON) simplify/__init__.py _examples/01deref-parameter/main.yaml > _examples/01deref-parameter/simplified.yaml
.PHONY: 01

validate-examples:
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/00*/main.yaml
	go run github.com/getkin/kin-openapi/cmd/validate@latest ./_examples/01*/main.yaml
.PHONY: validate-examples