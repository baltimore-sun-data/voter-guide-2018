#! /bin/bash
set -eux -o pipefail

# Get the directory that this script file is in
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd "$SCRIPT_DIR"

docker build -f dockerfiles/base.Dockerfile -t voter-guide-2018:base .
docker build -f dockerfiles/test.Dockerfile -t voter-guide-2018:tester .
docker run \
    -t \
    --rm \
    voter-guide-2018:tester
