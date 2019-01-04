#!/usr/bin/env bash

set -e
set -o pipefail

pushd $(dirname $(dirname $0)) > /dev/null # GO TO PROJECT ROOT

echo "Working dir: $(pwd -P)"

PLATFORM="$(uname -s | tr 'A-Z' 'a-z')"

export CGO_ENABLED=0
#GOOS=linux
export GOOS=${PLATFORM}

BUILD_DATE=$(date '+%Y%m%d%H%M%S%Z')
BUILD_LDFLAGS="-X github.com/argcv/gsp2md/version.GitHash=$(git rev-parse HEAD | cut -c1-8) "
BUILD_LDFLAGS="${BUILD_LDFLAGS} -X github.com/argcv/gsp2md/version.BuildDate=\"${BUILD_DATE}\" "
BUILD_LDFLAGS="${BUILD_LDFLAGS} \"-extldflags='-static'\""

function go-build() {
    build_path=$1
    echo "Fetching dependencies... ${build_path}"
    go get -d ${build_path}
    echo "Building... ${build_path}"
    go build -a -ldflags="$BUILD_LDFLAGS" ${build_path}
    echo "${build_path} is built"
    go test ${build_path} -cover
}

go-build ./cmd/gsp2md

unset BUILD_LDFLAGS
unset BUILD_DATE
unset GOOS
unset CGO_ENABLED
unset PLATFORM

popd > /dev/null # EXIT FROM PROJECT ROOT
