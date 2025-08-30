BINARY=sb
SRC=./cmd/sb.go
PREFIX=$(HOME)/.local/bin

# ============================================================
# BUILD
# ============================================================

.PHONY: build install test help confirm

## build: build the binary
build:
	go build -o $(BINARY) $(SRC)

## install: build the binary and install it to /.local/bin
install: build
	install -d $(PREFIX)
	install $(BINARY) $(PREFIX)/$(BINARY)

.PHONY: build-all
## build-all: build binaries for macOS, Linux, and Windows (amd64)
build-all: build
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/sb-darwin-amd64 ./cmd/sb.go
	GOOS=linux GOARCH=amd64 go build -o bin/sb-linux-amd64 ./cmd/sb.go
	GOOS=windows GOARCH=amd64 go build -o bin/sb-windows-amd64.exe ./cmd/sb.go

# ============================================================
# TESTS
# ============================================================		
	
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