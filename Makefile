PKG_LIST := $(shell go list ./...)

.PHONY: setup
setup:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.17.1

.PHONY: modules
modules:
	go mod tidy

.PHONY: build
build:
	go build

.PHONY: lint
lint:
	./bin/golangci-lint run

.PHONY: test
test: 
	go clean -testcache ${PKG_LIST}
	go test -short --race ${PKG_LIST}

default: build
