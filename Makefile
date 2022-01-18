make build:
	@go build

test:
	@go fmt
	@go test ./... -v

coverage:
	echo "Test Coverage script will be here"

clean:
	rm pilot-manager
