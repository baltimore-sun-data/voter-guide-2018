#! /bin/bash
set -eux -o pipefail

if [[ -n "$(which go)" ]]; then
    export GO111MODULE=on
    go run ./cmd/robocopy \
        -local \
        -metadata-src data/primary-metadata.json \
        -results-src data/primary-results.json \
        -output-dir dist/primary-results/
    go run ./cmd/robocopy \
        -local \
        -metadata-src data/general-metadata.json \
        -results-src data/general-results.json \
        -output-dir dist/results/
else
    robocopy \
        -local \
        -metadata-src data/primary-metadata.json \
        -results-src data/primary-results.json \
        -output-dir dist/primary-results/
    robocopy \
        -local \
        -metadata-src data/general-metadata.json \
        -results-src data/general-results.json \
        -output-dir dist/results/
fi
