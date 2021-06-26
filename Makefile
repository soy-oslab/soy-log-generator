GOBIN=go
GORUN=$(GOBIN) run
GOBUILD=$(GOBIN) build
GOTEST=$(GOBIN) test
BINPATH=$(shell pwd)/bin

all: generator

test: compressor-test

clean:
	rm -rf $(BINPATH)/*

generator:
	cd ./cmd/generator/; $(GORUN) .

generator-build:
	mkdir -p $(BINPATH)/generator
	cd ./cmd/generator/; $(GOBUILD) -o $(BINPATH)/generator

compressor-test:
	cd ./internal/compressor/; $(GOTEST) .

