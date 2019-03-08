NAME	:= myaws

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps
	go build -o bin/$(NAME)

.PHONY: package
package:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release:
	goreleaser --rm-dist

.PHONY: lint
lint:
	@golint ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test ./...

.PHONY: check
check: lint vet test build
