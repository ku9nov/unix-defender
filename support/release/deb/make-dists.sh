#!/bin/sh

set -eux

VERSION=$(cat unix-defender/DEBIAN/control |grep "^Version"|awk '{ print $2 }')
fakeroot dpkg-deb --build unix-defender output/unix-defender_$VERSION.deb
chown -R $TARGET_UID:$TARGET_GID output/