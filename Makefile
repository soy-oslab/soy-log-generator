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
	cd ./pkg/compressor; $(GOTEST) -cover -v -coverprofile=../../coverage.out .
	go tool cover -func=coverage.out
	rm coverage.out

compressor-bench:
	$(GOTEST) -bench=. -benchmem -benchtime=$(BENCHTIME) -timeout $(BENCHTIMEOUT) ./pkg/compressor

buffering-test:
	cd ./pkg/buffering; $(GOTEST) -cover -v -coverprofile=../../coverage.out .
	go tool cover -func=coverage.out
	rm coverage.out

watcher-test:
	cd ./pkg/watcher; $(GOTEST) -cover -v -coverprofile=../../coverage.out .
	go tool cover -func=coverage.out
	rm coverage.out
