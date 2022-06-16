BIN     := speechly
VERSION ?= latest
SRC     = $(shell find cmd -type f -name '*.go')
ifneq ("$(wildcard decoder/lib/libspeechly*)","")
TAGS=on_device
else
TAGS=
endif


all: build test lint

build: bin/speechly

bin/speechly: $(shell git ls-files)
	go build -ldflags="-X 'github.com/speechly/cli/cmd.version=$(VERSION)'" -tags "$(TAGS)" -o bin/speechly

test:
	go test -v ./...

clean:
	rm -rf bin/ dist/

lint:
	golangci-lint run --exclude-use-default=false

fmt:
	gofmt -l -w $(SRC)

.PHONY: all build lint clean fmt
