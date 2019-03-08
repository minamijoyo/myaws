NAME	:= myaws

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

LINT := $(GOBIN)/golint
GORELEASER := $(GOBIN)/goreleaser

$(LINT): ; @go get github.com/golang/lint/golint
$(GORELEASER): ; @go get github.com/goreleaser/goreleaser

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go get -d -v .

.PHONY: build
build: deps
	go build -o bin/$(NAME)

.PHONY: package
package: $(goreleaser)
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: $(goreleaser)
	goreleaser --rm-dist

.PHONY: lint
lint: $(LINT)
	@golint ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test ./...

.PHONY: check
check: lint vet test build
