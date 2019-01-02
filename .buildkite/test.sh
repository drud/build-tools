#!/bin/bash

# This script is used to build drud/build-tools using buildkite


set -o errexit
set -o pipefail
set -o nounset
set -x

BUILD_OS=$(go env GOOS)
echo "--- buildkite building $BUILDKITE_JOB_ID at $(date) on $HOSTNAME for OS=$(go env GOOS) in $PWD"

# Our testbot should now be sane, run the testbot checker to make sure.
echo "--- Checking for sane testbot"
./.buildkite/sanetestbot.sh

echo "--- make $BUILD_OS"
cd tests
time make $BUILD_OS
echo "--- make test"
time make test
RV=$?
echo "--- build.sh completed with status=$RV"
exit $RV
