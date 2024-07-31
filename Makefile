UBINDIR ?= /usr/bin
DESTDIR ?=
EXTNAME := $(shell basename $(shell pwd))

# go tool nm ./anise | grep Commit
override LDFLAGS += -X "github.com/macaroni-os/anise-portage-converter/pkg/converter.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
override LDFLAGS += -X "github.com/macaroni-os/anise-portage-converter/pkg/converter.BuildCommit=$(shell git rev-parse HEAD)"

all: build install

build:
	CGO_ENABLED=0 go build -o anise-portage-converter -ldflags '$(LDFLAGS)'

install: build
	install -d $(DESTDIR)/$(UBINDIR)
	install -m 0755 $(EXTNAME) $(DESTDIR)/$(UBINDIR)/

.PHONY: deps
deps:
	go env
	# Installing dependencies...
	GO111MODULE=on go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
	go get github.com/onsi/gomega/...
	ginkgo version

.PHONY: test
test: deps
	ginkgo -r -race -flake-attempts 3 ./...

.PHONY: coverage
coverage:
	go test ./... -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: test-coverage
test-coverage:
	scripts/ginkgo.coverage.sh --codecov

.PHONY: clean
clean:
	-rm anise-portage-converter
	-rm -rf release/ dist/

.PHONY: goreleaser-snapshot
goreleaser-snapshot:
	rm -rf dist/ || true
	goreleaser release --skip=validate,publish --snapshot --verbose

