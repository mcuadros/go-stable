# Package configuration
PROJECT = go-stable
COMMANDS = cli/stable
DEPENDENCIES = golang.org/x/tools/cmd/cover
PACKAGES = .
ORGANIZATION = mcuadros

# Environment
BASE_PATH := $(shell pwd)
BUILD_PATH := $(BASE_PATH)/build
BUILD ?= $(shell date +"%m-%d-%Y_%H_%M_%S")
COMMIT ?= $(shell git log --format='%H' -n 1 | cut -c1-10)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

# Packages content
PKG_OS = darwin linux
PKG_ARCH = amd64
PKG_CONTENT = LICENSE
PKG_TAG = latest

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOGET = $(GOCMD) get -v
GOTEST = $(GOCMD) test -v
GHRELEASE = github-release

# Coverage
COVERAGE_REPORT = coverage.txt
COVERAGE_MODE = atomic

# Docker
DOCKERCMD = docker

ifneq ($(origin TRAVIS_TAG), undefined)
	BRANCH := $(TRAVIS_TAG)
endif

# Rules
all: clean upload

dependencies:
	$(GOGET) -t ./...
	for i in $(DEPENDENCIES); do $(GOGET) $$i; done

test: dependencies
	for p in $(PACKAGES); do \
		$(GOTEST) $${p}; \
	done;

test-coverage: dependencies
	echo "mode: $(COVERAGE_MODE)" > $(COVERAGE_REPORT); \
	for p in $(PACKAGES); do \
		$(GOTEST) $${p} -coverprofile=tmp_$(COVERAGE_REPORT) -covermode=$(COVERAGE_MODE); \
		cat tmp_$(COVERAGE_REPORT) | grep -v "mode: $(COVERAGE_MODE)" >> $(COVERAGE_REPORT); \
		rm tmp_$(COVERAGE_REPORT); \
	done;

packages: dependencies
	for os in $(PKG_OS); do \
		for arch in $(PKG_ARCH); do \
			cd $(BASE_PATH); \
			mkdir -p $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}; \
			for cmd in $(COMMANDS); do \
				cd $${cmd} && GOOS=$${os} GOARCH=$${arch} $(GOCMD) build -ldflags "-X main.version=$(BRANCH) -X main.build=$(BUILD) -X main.commit=$(COMMIT)" -o $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}/$$(basename $${cmd}) .; \
				cd $(BASE_PATH); \
			done; \
			for content in $(PKG_CONTENT); do \
				cp -rf $${content} $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}/; \
			done; \
			cd  $(BUILD_PATH) && tar -cvzf $(BUILD_PATH)/$(PROJECT)_$${os}_$${arch}.tar.gz $(PROJECT)_$${os}_$${arch}/; \
		done; \
	done;

push:
	for cmd in $(COMMANDS); do \
		cd $${cmd} && CGO_ENABLED=0 $(GOCMD) build -ldflags "-X main.version=$(BRANCH) -X main.build=$(BUILD) -X main.commit=$(COMMIT)" .; \
		cd $(BASE_PATH); \
	done;
	$(DOCKERCMD) build -t $(ORGANIZATION)/$(PROJECT):$(BRANCH) .
	$(DOCKERCMD) push $(ORGANIZATION)/$(PROJECT):$(BRANCH)

clean:
	rm -rf $(BUILD_PATH)
	$(GOCLEAN) .
