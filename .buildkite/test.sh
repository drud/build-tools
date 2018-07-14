#!/bin/bash

# This script is used to build drud/build-tools using buildkite

# Manufacture a $GOPATH environment that can mount on docker (when buildkite build)
export GOPATH=~/tmp/buildkite-fake-gopath/$BUILDKITE_JOB_ID
DRUDSRC=$GOPATH/src/github.com/drud
mkdir -p $DRUDSRC
ln -s $PWD $DRUDSRC/build-tools
cd $DRUDSRC/build-tools
BUILD_OS=$(go env GOOS)
echo "--- buildkite building $BUILDKITE_JOB_ID at $(date) on $HOSTNAME for OS=$(go env GOOS) in $PWD GOPATH=$GOPATH"

echo "--- cleaning up docker"
echo "Warning: deleting all docker containers"
if [ "$(docker ps -aq | wc -l)" -gt 0 ] ; then
	docker rm -f $(docker ps -aq)
fi
docker system prune --force

# Update all images that could have changed
docker images | awk '/drud/ {print $1":"$2 }' | xargs -L1 docker pull

set -o errexit
set -o pipefail
set -o nounset
set -x

# Our testbot should now be sane, run the testbot checker to make sure.
./.buildkite/sanetestbot.sh

echo "--- make $BUILD_OS"
cd tests
time make $BUILD_OS
echo "--- make test"
time make test
RV=$?
echo "--- build.sh completed with status=$RV"
exit $RV
