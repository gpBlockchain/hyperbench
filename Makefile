# Go parameters
GOCMD=go
PACKR=packr
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGEN=$(GOCMD) generate
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
BINARY_NAME=hyperbench
SKIP_DIR=benchmark

# version info
branch=$(shell git rev-parse --abbrev-ref HEAD)
commitID=$(shell git log --pretty=format:"%h" -1)
date=$(shell date +%Y%m%d)
importpath=github.com/meshplus/hyperbench/cmd
ldflags=-X ${importpath}.branch=${branch} -X ${importpath}.commitID=${commitID} -X ${importpath}.date=${date}

# path
ASSETS=filesystem/assets
DIRS=benchmark
GET=github.com/gobuffalo/packr/v2/... github.com/gobuffalo/packr/v2/packr2
FAILPOINT=github.com/pingcap/failpoint/failpoint-ctl

all: build

## build: build the binary with pre-packed static resource
build: dep assets
	@export GOPROXY=https://goproxy.cn,direct
	@packr2 build -o $(BINARY_NAME) -ldflags "${ldflags}"
	@-rm -rf $(ASSETS)

## pack: build the binary with local static resource
pack: assets
	@packr2 build -o $(BINARY_NAME) -ldflags "${ldflags}"
	@-rm -rf $(ASSETS)

## test: run all test
test:
	@go get $(FAILPOINT)
	@failpoint-ctl enable
	@$(GOTEST) `go list ./... | grep -v $(SKIP_DIR)`
	@failpoint-ctl disable

## clean: clean all file generated by make
clean:
	@packr2 clean
	@-rm -rf $(BINARY_NAME)
	@-rm -rf $(ASSETS)

.PHONY: assets
## assets: prepare asserts
assets:
	@-rm -rf $(ASSETS)
	@mkdir $(ASSETS)
	@cp -r $(DIRS) $(ASSETS)

.PHONY: dep
## dep: install the dependencies outside (may need to use proxy to download some packages)
dep:
	@go get -u $(GET)

help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
