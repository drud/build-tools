# Makefile for a standard repo with associated container

##### These variables need to be adjusted in most repositories #####

# This repo's root import path (under GOPATH).
PKG := github.com/drud/build-tools/tests

# Docker repo for a push
DOCKER_REPO ?= drud/build-tools-test

# Upstream repo used in the Dockerfile
UPSTREAM_REPO ?= golang:1.7.5-alpine3.5

export WORKING_DIR = $(shell pwd)

# Top-level directories to build - Only build the explicit complex fail stuff by default. Unit test will be overridden.
SRC_DIRS := cmd pkg

# Version variables to replace in build, The variable VERSION is automatically pulled from git committish so it doesn't have to be added
# These are replaced in the $(PKG).version package DBImg DBTag RouterImage RouterTag
#VERSION_VARIABLES = DdevVersion WebImg

# These variables will be used as the default unless overridden by the make command line
#DdevVersion ?= $(VERSION)
#WebImg ?= drud/nginx-php-fpm7-local

# Optional to docker build
# DOCKER_ARGS =

# VERSION can be set by
  # Default: git tag
  # make command line: make VERSION=0.9.0
# It can also be explicitly set in the Makefile as commented out below.

# This version-strategy uses git tags to set the version string
# VERSION can be overridden on make commandline: make VERSION=0.9.1 push
VERSION := $(shell git describe --tags --always --dirty)
#
# This version-strategy uses a manual value to set the version string
#VERSION := 1.2.3

# Each section of the Makefile is included from standard components below.
# If you need to override one, import its contents below and comment out the
# include. That way the base components can easily be updated as our general needs
# change.
include ../makefile_components/base_build_go.mak
#include ../makefile_components/base_build_python-docker.mak
include ../makefile_components/base_container.mak
include ../makefile_components/base_push.mak
#include ../makefile_components/base_test_go.mak
#include ../makefile_components/base_test_python.mak


# We can't use the standard 'test' target because this one actually uses 'make' and counts on resources unavailable in
# golang compiler container.
test: build
	@mkdir -p bin/linux
	@mkdir -p  $(GOTMP)/{src/$(PKG),pkg,bin,std/linux}
	go test -v -installsuffix "static" -ldflags '$(LDFLAGS)' $(SRC_AND_UNDER) $(TESTARGS)
	$(MAKE) -C standard_target $@
	$(MAKE) -C with_vendor $@


# test_precompile allows a full compilation of _test.go files, without execution of the tests.
# Setup and teardown in TestMain is still executed though, so this can cost some time.
test_precompile: TESTARGS=-run '^$$'
test_precompile: test

# Simple way to execute a random command in the container for tests - used only for testing
# Example: make COMMAND="govendor fetch golang.org/x/net/context"
COMMAND := govendor list
container_cmd:
	@docker run                                                            \
		-t --rm -u $(shell id -u):$(shell id -g)                         \
	    -v $(PWD)/$(GOTMP):/go                                                 \
	    -v $(PWD):/go/src/$(PKG)                                          \
	    -v $(PWD)/bin/linux:/go/bin                                     \
	    -v $(PWD)/$(GOTMP)/std/linux:/usr/local/go/pkg/linux_amd64_static  \
		-w /go/src/$(PKG)                                                  \
		$(BUILD_IMAGE)                                                     \
		$(COMMAND)

