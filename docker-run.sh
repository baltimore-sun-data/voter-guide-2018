#! /bin/bash
set -eux -o pipefail

# Get the directory that this script file is in
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
cd "$SCRIPT_DIR"

docker build . -t voter-guide-2018-deployer
docker run \
    -t \
    --rm \
    -e AWS_ACCESS_KEY_ID \
    -e AWS_SECRET_ACCESS_KEY \
    voter-guide-2018-deployer
