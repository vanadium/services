SHELL := /bin/bash -euo pipefail

VDLPATH ?= $(shell pwd)
export VDLPATH
vdlgen:
	go run v.io/x/ref/cmd/vdl generate --show-warnings=true --lang=go v.io/...
	go generate ./v23/vdl/vdltest

.PHONY: test
test:
	@echo "VDLPATH" "${VDLPATH}"
	go test ./...

.PHONY: test-integration
test-integration:
	@echo "VDLPATH" "${VDLPATH}"
	go test \
		github.com/vanadium/services/identity/identityd \
		-v23.tests

refresh:
	go generate ./...
	go get cloudeng.io/go/cmd/goannotate
	go run cloudeng.io/go/cmd/goannotate --config=vanadium-code-annotations.yaml --annotation=copyright ./...
	go mod tidy
