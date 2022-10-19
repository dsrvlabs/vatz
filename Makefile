VERSION := $(shell git describe --tags)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-X 'main.Version=$(BRANCH)' -X 'main.Commit=$(COMMIT_HASH)'"

.PHONY: test build coverage clean

all: test coverage build

test:
	@go fmt
	@go test ./... -v

coverage:
	echo "Test Coverage script will be here"
	@go test -coverprofile cover.out ./...

build:
	go build $(LDFLAGS) -v

clean:
	go clean
