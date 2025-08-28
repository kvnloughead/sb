BINARY=sb
SRC=./cmd/sb.go
PREFIX=$(HOME)/.local/bin


.PHONY: build install build-defaults test help confirm
## build-defaults: generate Go source from defaults.yaml
build-defaults:
	go run scripts/gen_defaults.go defaults.yaml config/generated_defaults.go

## build: build the binary
build: build-defaults
	go build -o $(BINARY) $(SRC)

## install: build the binary and install it to /.local/bin
install: build
	install -d $(PREFIX)
	install $(BINARY) $(PREFIX)/$(BINARY)
	rm $(BINARY)
	
## test: run Go tests
test:
	go test ./...

# ============================================================
# HELPERS
# ============================================================

## help: print this help message
.PHONY: help
help:
	@echo "\nUsage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
	@echo

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]