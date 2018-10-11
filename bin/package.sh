#!/bin/bash

set -e

VERSION=$1
SOURCE_DIR=$2
TARGET_PATH=$3

if [ -z "${SOURCE_DIR}" ] || [ -z "${VERSION}" ]; then
  echo "You must give source_dir and version as arguments!"
  echo "Usage: package.sh <version number> <source_dir> <target package path>"
  exit 1;
fi

if [ ! `hash fpm 2>/dev/null` ]; then
  echo "Installing Ruby and friends for FPM"
  apt-get install -y ruby-dev gcc

  echo "Installing FPM"
  gem install fpm
fi

echo "Packaging with FPM"

# Build the deb -package
fpm -s dir -t deb \
  --name vanadia \
  -v $VERSION \
  -C $SOURCE_DIR \
  --package $TARGET_PATH \
  --description "API Blueprint to Postman collection converter" \
  .