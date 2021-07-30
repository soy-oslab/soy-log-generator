GOBIN=go
GORUN=$(GOBIN) run
GOBUILD=$(GOBIN) build
GOTEST=$(GOBIN) test


BENCHTIME=1s
BENCHTIMEOUT=10m

ifeq ($(OS),Windows_NT)
	MKDIR=mkdir
	RM=del
	RMDIR=rmdir /S /Q
	RMFLAG=/F /Q
	SEP=\\
else
	MKDIR=mkdir -p
	RM=rm
	RMDIR=rmdir
	RMFLAG=-rf
	SEP=/
endif

BUILD_PATH=.$(SEP)build

all: generator-build

test: compressor-test buffering-test watcher-test \
      scheduler-test ring-test transport-test \
      classifier-test

clean:
ifeq ($(OS), Windows_NT)
	if exist "$(BUILD_PATH)" $(RM) $(RMFLAG) $(BUILD_PATH)$(SEP)*
	if exist "$(BUILD_PATH)" $(RMDIR) $(BUILD_PATH)
else
	$(RM) $(RMFLAG) $(BUILD_PATH)$(SEP)*
	$(RMDIR) $(BUILD_PATH)
endif

generator-run:
	$(GORUN) .$(SEP)cmd$(SEP)generator$(SEP)main.go

generator-build:
ifeq ($(OS), Windows_NT)
	if not exist "$(BUILD_PATH)$(SEP)generator" $(MKDIR) $(BUILD_PATH)$(SEP)generator
else
	$(MKDIR) $(BUILD_PATH)$(SEP)generator
endif
	$(GOBUILD) -o $(BUILD_PATH)$(SEP)generator .$(SEP)cmd$(SEP)generator$(SEP)main.go

compressor-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)compressor
	go tool cover -func=coverage.out
	$(RM) coverage.out

compressor-bench:
	$(GOTEST) -bench=. -benchmem -benchtime=$(BENCHTIME) -timeout $(BENCHTIMEOUT) .$(SEP)pkg$(SEP)compressor

buffering-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)buffering
	go tool cover -func=coverage.out
	$(RM) coverage.out

watcher-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)watcher
	go tool cover -func=coverage.out
	$(RM) coverage.out

scheduler-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)scheduler
	go tool cover -func=coverage.out
	$(RM) coverage.out

ring-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)ring
	go tool cover -func=coverage.out
	$(RM) coverage.out

transport-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)transport
	go tool cover -func=coverage.out
	$(RM) coverage.out

classifier-test:
	$(GOTEST) -cover -v -coverprofile=coverage.out .$(SEP)pkg$(SEP)scheduler
	go tool cover -func=coverage.out
	$(RM) coverage.out

codacy-coverage-push:
ifeq ($(OS),Windows_NT)
	@echo "Windows_NT doesn't support"
else
	$(GOTEST) -coverprofile=coverage.out .$(SEP)...
	bash scripts$(SEP)get.sh report --force-coverage-parser go -r .$(SEP)coverage.out
endif

.PHONY: gen-src-archive
gen-src-archive:
ifeq ($(OS),Windows_NT)
	@echo "Windows_NT doesn't support"
else
	bash scripts$(SEP)gen_src_archive.sh
endif
