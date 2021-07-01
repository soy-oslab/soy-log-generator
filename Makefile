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
	$(GORUN) ./cmd/generator/main.go

generator-build:
	mkdir -p $(BINPATH)/generator
	$(GOBUILD) -o $(BINPATH)/generator ./cmd/generator/main.go

compressor-test:
	$(GOTEST) -cover -v ./pkg/compressor

compressor-bench:
	$(GOTEST) -bench=. -benchmem -benchtime=$(BENCHTIME) -timeout $(BENCHTIMEOUT) ./pkg/compressor

