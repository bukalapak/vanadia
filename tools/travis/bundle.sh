#!/bin/bash

TRAVIS_ARCHIVE_FILE=vanadia-$TRAVIS_TAG.${TRAVIS_OS_NAME/osx/darwin}-amd64.tar.gz

if [[ "$TRAVIS_TEST_RESULT" == "0" ]]; then
    echo "Generating ${TRAVIS_ARCHIVE_FILE}"
    tar -czf $TRAVIS_ARCHIVE_FILE vanadia LICENSE README.md vanadia.yml
    shasum -a 256 -b $TRAVIS_ARCHIVE_FILE
fi
