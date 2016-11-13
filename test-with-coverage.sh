#!/bin/bash

set -e -u
set -o pipefail

ginkgo -r -p --randomizeAllSpecs --randomizeSuites --failOnPending \
    --trace --race --compilers=2 \
    -coverpkg=github.com/rosenhouse/umbrella/... \
    -covermode=atomic -coverprofile=umbrella.coverprofile

# find . -name "*.coverprofile*" | xargs gocovmerge > all.coverprofile
