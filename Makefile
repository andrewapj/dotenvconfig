.PHONY: all clean fmt vet test 

# Default target
all: test

clean:
	go clean .
	rm -f coverage.html coverage.out

fmt: clean
	go fmt ./...

vet: fmt
	go vet ./...

test: vet
	go test -v -coverprofile coverage.out ./... && go tool cover -html coverage.out -o coverage.html
