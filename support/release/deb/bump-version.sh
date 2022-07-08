#!/usr/bin/env bash

set -x
export VERSION=$(cat VERSION)
export OLDVER=${VERSION}
export MAJOR=$(echo ${VERSION} | cut -d'.' -f1)
export MINOR=$(echo ${VERSION} | cut -d'.' -f2)
export MINOR=$((MINOR+1))
export VERSION="${MAJOR}.${MINOR}"
echo ${VERSION}>VERSION
sed -i "s/${OLDVER}/${VERSION}/" unix-defender/DEBIAN/control