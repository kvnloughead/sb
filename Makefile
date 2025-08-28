BINARY=sb
SRC=./cmd/sb.go
PREFIX=$(HOME)/.local/bin

.PHONY: build install

## build: build the binary
build:
	go build -o $(BINARY) $(SRC)

## install: build the binary and install it to /.local/bin
install: build
	install -d $(PREFIX)
	install $(BINARY) $(PREFIX)/$(BINARY)
	rm $(BINARY)

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