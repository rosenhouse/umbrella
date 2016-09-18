#!/bin/bash

set -e -u

ginkgo -r -p --randomizeAllSpecs --randomizeSuites --failOnPending \
    --trace --race --compilers=2 \
    -coverpkg=github.com/rosenhouse/umbrella/... \
    -covermode=set -coverprofile=umbrella.coverprofile

set +e
find . -name "*.coverprofile*" | xargs gocovmerge > all.coverprofile
