BIN := speechly

build: bin/speechly

bin/speechly: $(shell git ls-files)
	go build -o bin/speechly

clean:
	rm -rf bin/ dist/

lint:
	golangci-lint run --exclude-use-default=false

.PHONY: build lint clean
