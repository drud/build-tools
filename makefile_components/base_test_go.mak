# test section of Makefile
##### PLEASE DO NOT CHANGE THIS FILE #####
##### If one of these sections does not meet your needs, consider copying its
##### contents into ../Makefile and commenting out the include and adding a
##### comment about what you did and why.

TESTOS = $(BUILD_OS)

test: build
	@echo "Testing $(SRC_AND_UNDER) with TESTARGS=$(TESTARGS)"
	@mkdir -p $(GOTMP)/{.cache,pkg,src,bin}
	docker run -t --rm  -u $(shell id -u):$(shell id -g)                 \
	    -v $(PWD):/workdir           \
	    -v $(PWD)/$(GOTMP)/bin:/go/bin                 \
	    -e CGO_ENABLED=0	\
	    -e GOPATH=/workdir/.gotmp \
	    -w /workdir                                                  \
	    $(BUILD_IMAGE)                                                     \
        go test $(USEMODVENDOR) -v -installsuffix static -ldflags '$(LDFLAGS)' $(SRC_AND_UNDER) $(TESTARGS)
	@$(shell chmod -R u+w $(GOTMP))

# test_precompile allows a full compilation of _test.go files, without execution of the tests.
# Setup and teardown in TestMain is still executed though, so this can cost some time.
test_precompile: TESTARGS=-run '^$$'
test_precompile: test
