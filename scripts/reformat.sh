#!/usr/bin/env bash

set -e
set -o pipefail

pushd $(dirname $(dirname $0)) > /dev/null # GO TO PROJECT ROOT

echo "Working dir: $(pwd -P)"

SRC_DIRS=('configs' 'pkg')

for SRC_DIR in ${SRC_DIRS[@]}; do
    echo "Processing folder: ${SRC_DIR}"
    find ${SRC_DIR} -name '*.go' | xargs -n 1 -I{} -P 6 sh -c 'echo "reformat: {}" && gofmt -w {}'
done



popd > /dev/null # EXIT FROM PROJECT ROOT
