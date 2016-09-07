#!/bin/bash

set -e
set -x

VERSION=$(grep "const version " cmd/version.go | sed -E 's/.*"(.+)"$/\1/')
REPO="myaws"

rm -rf ./out/
gox --osarch "darwin/amd64 linux/amd64" -output="./out/${REPO}_${VERSION}_{{.OS}}_{{.Arch}}/{{.Dir}}"

rm -rf ./pkg/
mkdir ./pkg

for PLATFORM in $(find ./out -mindepth 1 -maxdepth 1 -type d); do
    PLATFORM_NAME=$(basename ${PLATFORM})

    pushd ${PLATFORM}
    zip ../../pkg/${PLATFORM_NAME}.zip ./*
    popd
done
