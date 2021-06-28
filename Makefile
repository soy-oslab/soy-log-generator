GOBIN=go
GORUN=$(GOBIN) run
GOBUILD=$(GOBIN) build
GOTEST=$(GOBIN) test

BINPATH=$(shell pwd)/bin

BENCHTIME=1s
BENCHTIMEOUT=10m

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
	cd ./internal/compressor/; $(GOTEST) -v -bench=. -benchmem -benchtime=$(BENCHTIME) -timeout $(BENCHTIMEOUT)

