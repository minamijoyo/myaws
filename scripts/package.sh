#!/bin/bash

set -e
set -x

VERSION="0.0.1"
REPO="myaws"

rm -rf ./out/
gox --osarch "darwin/amd64 linux/amd64" -output="./out/${REPO}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

rm -rf ./pkg/
mkdir ./pkg

for PLATFORM in $(find ./out -mindepth 1 -maxdepth 1 -type d); do
    PLATFORM_NAME=$(basename ${PLATFORM})
    ARCHIVE_NAME=${REPO}_${VERSION}_${PLATFORM_NAME}

    pushd ${PLATFORM}
    zip ../../pkg/${ARCHIVE_NAME}.zip ./*
    popd
done

