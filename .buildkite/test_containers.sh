#!/bin/bash

# Runs "make test" in each container directory

set -o errexit
set -o pipefail
set -o nounset
set -x
echo "--- buildkite building ${BUILDKITE_JOB_ID:-} at $(date) on $(hostname) for OS=$(go env GOOS) in ${DRUDSRC:-}/ddev"

echo "--- Cleanup docker"
echo "Warning: deleting all docker containers and deleting images that match this build."
if [ "$(docker ps -aq | wc -l)" -gt 0 ] ; then
	docker rm -f $(docker ps -aq)
fi
docker system prune --volumes --force

# Make sure we don't have any existing containers on the testbot that might
# result in this container not being built from scratch.
VERSION=$(make version | sed 's/^VERSION://')
IMAGES=$(docker images | awk "/$VERSION/ { print \$3 }")
if [ ! -z "$IMAGES" ] ; then
  docker rmi --force $IMAGES
fi

# Our testbot should now be sane, run the testbot checker to make sure.
./.buildkite/sanetestbot.sh

for dir in containers/*
    do pushd $dir
    echo "--- Build container $dir"
    time make container
    echo "--- Test container $dir"
    time make test
    popd
done
