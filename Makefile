BIN := speechly

SRC = $(shell find cmd -type f -name '*.go')

build: bin/speechly

bin/speechly: $(shell git ls-files)
	go build -o bin/speechly

test:
	go test -v ./...

clean:
	rm -rf bin/ dist/

lint:
	golangci-lint run --exclude-use-default=false

fmt:
	gofmt -l -w $(SRC)

.PHONY: build lint clean fmt
