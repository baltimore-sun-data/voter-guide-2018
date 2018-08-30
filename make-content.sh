#! /bin/bash
set -eux -o pipefail

if [[ -n "$(which go)" ]]; then
    go run cmd/csv-splitter/*.go csv/*.csv
else
    csv-splitter csv/*.csv
fi
