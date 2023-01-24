all: clean build test lint

clean:
	rm -rf bin/op-rosetta

build:
	env GO111MODULE=on go build -o bin/op-rosetta ./cmd

test:
	go test -v ./...

lint:
	golangci-lint run -E asciicheck,goimports,misspell ./...

.PHONY: \
	test \
	lint \
	build
