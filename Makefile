GOBIN=go
GORUN=$(GOBIN) run
GOBUILD=$(GOBIN) build
GOTEST=$(GOBIN) test

BUILD_PATH=$(shell pwd)/build

BENCHTIME=1s
BENCHTIMEOUT=10m

all: generator-build

test: compressor-test buffering-test watcher-test

clean:
	rm -rf $(BUILD_PATH)/*

generator-run:
	$(GORUN) ./cmd/generator/main.go

generator-build:
	mkdir -p $(BUILD_PATH)/generator
	$(GOBUILD) -o $(BUILD_PATH)/generator ./cmd/generator/main.go

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

codacy-coverage-push:
	$(GOTEST) -coverprofile=coverage.out ./...
	bash scripts/get.sh report --force-coverage-parser go -r ./coverage.out

.PHONY: gen-src-archive
gen-src-archive:
	bash scripts/gen_src_archive.sh
