#!/usr/bin/env bash

set -x
export VERSION=$(cat VERSION)
export OLDVER=${VERSION}
export MAJOR=$(echo ${VERSION} | cut -d'.' -f1)
export MINOR=$(echo ${VERSION} | cut -d'.' -f2)
export PATCH=$(echo ${VERSION} | cut -d'.' -f3)
export PATCH=$((PATCH+1))
export VERSION="${MAJOR}.${MINOR}.${PATCH}"
echo ${VERSION}>VERSION
if [[ "$OSTYPE" == "darwin"* ]]; then
    gsed -i "s/${OLDVER}/${VERSION}/" unix-defender/DEBIAN/control
else
    sed -i "s/${OLDVER}/${VERSION}/" unix-defender/DEBIAN/control
fi
