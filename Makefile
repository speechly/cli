BIN     := speechly
VERSION ?= latest
SRC     = $(shell find cmd -type f -name '*.go')
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	PLATFORM=macos
else ifeq ($(UNAME_S),Linux)
	PLATFORM=linux
endif
ifneq ("$(wildcard decoder/${PLATFORM}-x86_64/lib/libspeechly*)","")
TAGS=on_device
else
TAGS=
endif


all: build test lint docs

build: bin/speechly

bin/speechly: $(shell git ls-files)
	go build -ldflags="-X 'github.com/speechly/cli/cmd.version=$(VERSION)'" -tags "$(TAGS)" -o bin/speechly

tflite-version: TAGS += tflite
tflite-version: all

coreml-version: TAGS += coreml
coreml-version: all

test:
	go test -v ./...

docs:
	go run docs/generate.go docs

clean:
	rm -rf bin/ dist/

lint:
	golangci-lint run --exclude-use-default=false

fmt:
	gofmt -l -w $(SRC)

.PHONY: all build lint clean fmt docs tflite-version coreml-version
