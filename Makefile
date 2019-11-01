NAME	:= myaws

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

GOLINT := $(GOBIN)/golint
INEFFASSIGN := $(GOBIN)/ineffassign
GORELEASER := $(GOBIN)/goreleaser

$(GOLINT): ; @go install github.com/golang/lint/golint
$(INEFFASSIGN): ; @go install github.com/gordonklaus/ineffassign
$(GORELEASER): ; @go install github.com/goreleaser/goreleaser

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go mod download

.PHONY: build
build: deps
	go build -o bin/$(NAME)

.PHONY: package
package: $(GORELEASER)
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: release
release: $(GORELEASER)
	goreleaser --rm-dist

.PHONY: lint
lint: $(GOLINT)
	@golint ./...

.PHONY: ineffassign
ineffassign: $(INEFFASSIGN)
	@ineffassign ./

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test ./...

.PHONY: check
check: lint ineffassign vet test build
