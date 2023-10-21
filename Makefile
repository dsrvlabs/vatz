VERSION := $(shell git describe --tags)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-X 'github.com/dsrvlabs/vatz/utils.Version=$(BRANCH)' -X 'github.com/dsrvlabs/vatz/utils.Commit=$(COMMIT_HASH)'"

.PHONY: test build coverage clean lint

all: test coverage build

test:
	@go fmt
	@go test ./... -v

coverage:
	@go test -coverprofile cover.out ./...

build:
	go build $(LDFLAGS) -v

clean:
	go clean

lint:
	golangci-lint run --timeout 5m	
