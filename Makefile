PROTOFILES := $(shell find protos -name "*.proto")
PROTOGOFILES := $(subst protos/,gen/go/,$(PROTOFILES:.proto=.pb.go))
BIN := speechly

build: bin/speechly

bin/speechly: $(shell git ls-files)
bin/speechly: $(PROTOGOFILES)
	go build -o bin/speechly

$(PROTOGOFILES): gen/go/%.pb.go: protos/%.proto
	@mkdir -p gen/go
	@protoc -I protos --go_out=plugins=grpc:gen/go $<

clean:
	rm -rf bin/ dist/ gen/

lint:
	golangci-lint run --exclude-use-default=false

.PHONY: build lint clean
