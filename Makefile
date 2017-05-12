NAME				:= myaws
VERSION			:= v0.2.0
REVISION		:= $(shell git rev-parse --short HEAD)
LDFLAGS			:= "-X github.com/minamijoyo/myaws/cmd.Version=${VERSION} -X github.com/minamijoyo/myaws/cmd.Revision=${REVISION} -extldflags \"-static\""
OSARCH			:= "darwin/amd64 linux/amd64"
GITHUB_ORG	:= minamijoyo

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

LINT := $(GOBIN)/golint
GOX := $(GOBIN)/gox
ARCHIVER := $(GOBIN)/archiver
GHR := $(GOBIN)/ghr

$(LINT): ; @go get github.com/golang/lint/golint
$(GOX): ; @go get github.com/mitchellh/gox
$(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/archiver
$(GHR): ; @go get github.com/tcnksm/ghr

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go get -d -v .

.PHONY: build
build: deps
	go build -a -tags netgo -installsuffix netgo -ldflags $(LDFLAGS) -o bin/$(NAME)

.PHONY: install
install: deps
	go install -ldflags $(LDFLAGS)

.PHONY: cross-build
cross-build: deps $(GOX)
	rm -rf ./out && \
	gox -ldflags $(LDFLAGS) -osarch $(OSARCH) -output "./out/${NAME}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

.PHONY: package
package: cross-build $(ARCHIVER)
	rm -rf ./pkg && mkdir ./pkg && \
	pushd out && \
	find * -type d -exec archiver make ../pkg/{}.tar.gz {}/$(NAME) \; && \
	popd

.PHONY: release
release: $(GHR)
	ghr -u $(GITHUB_ORG) $(VERSION) pkg/

.PHONY: digest
digest:
	openssl dgst -sha256 pkg/${NAME}_${VERSION}_darwin_amd64.tar.gz

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
