CC=gcc

GOBIN=go
GORUN=$(GOBIN) run
GOBUILD=$(GOBIN) build
GOTEST=$(GOBIN) test


BENCHTIME=1s
BENCHTIMEOUT=10m

BUILD_PATH=./build

all: generator-build

test: compressor-test buffering-test watcher-test \
      scheduler-test ring-test transport-test \
      classifier-test

clean:
	rm $(RMFLAG) $(BUILD_PATH)/*
	rmdir $(BUILD_PATH)

generator-run:
	$(GORUN) ./cmd/generator/main.go

generator-build:
	mkdir -p $(BUILD_PATH)/generator
	$(GOBUILD) -o $(BUILD_PATH)/generator ./cmd/generator/main.go

kube-wrapper-build:
	mkdir -p $(BUILD_PATH)/kube
	$(CC) -Wall -Werror -o $(BUILD_PATH)/kube/wrapper ./tools/kube-generator-wrapper/kube-generator-wrapper.c

compressor-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/compressor
	go tool cover -func=coverage.out
	rm coverage.out

compressor-bench:
	$(GOTEST) -bench=. -benchmem -benchtime=$(BENCHTIME) -timeout $(BENCHTIMEOUT) ./pkg/compressor

buffering-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/buffering
	go tool cover -func=coverage.out
	rm coverage.out

watcher-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/watcher
	go tool cover -func=coverage.out
	rm coverage.out

scheduler-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/scheduler
	go tool cover -func=coverage.out
	rm coverage.out

ring-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/ring
	go tool cover -func=coverage.out
	rm coverage.out

transport-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/transport
	go tool cover -func=coverage.out
	rm coverage.out

classifier-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out ./pkg/scheduler
	go tool cover -func=coverage.out
	rm coverage.out

codacy-coverage-push:
	$(GOTEST) -coverprofile=coverage.out ./...
	bash scripts/get.sh report --force-coverage-parser go -r ./coverage.out

.PHONY: gen-src-archive
gen-src-archive:
	bash scripts/gen_src_archive.sh
