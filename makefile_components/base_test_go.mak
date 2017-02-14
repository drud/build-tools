# test section of Makefile
##### PLEASE DO NOT CHANGE THIS FILE #####
##### If one of these sections does not meet your needs, consider copying its
##### contents into ../Makefile and commenting out the include and adding a
##### comment about what you did and why.

TESTOS = $(shell uname -s | tr '[:upper:]' '[:lower:]')

test: linux
	@mkdir -p bin/linux
	@mkdir -p .go/src/$(PKG) .go/pkg .go/bin .go/std/linux
	@docker run                                                            \
	    -t                                                                \
	    -u root:root                                             \
		-v $(BUILD_BASE_DIR)/build-tools:/build-tools		\
	    -v $$(pwd)/.go:/go                                                 \
	    -v $$(pwd):/go/src/$(PKG)                                          \
	    -v $$(pwd)/bin/linux:/go/bin                                     \
	    -v $$(pwd)/.go/std/linux:/usr/local/go/pkg/linux_amd64_static  \
	    -w /go/src/$(PKG)                                                  \
	    -e GOOS=linux	\
	    $(BUILD_IMAGE)                                                     \
	    /bin/sh -c "                                                       \
	        GOOS=$(TESTOS)                                                   \
	        go test -v -installsuffix "static" ${TARGETS} $(SRC_AND_BELOW)   \
	    "
