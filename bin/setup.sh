#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")/.."

# install govendor
go get -u github.com/kardianos/govendor

# install govendor
govendor sync

# build
apt-get update && apt-get install -y make && make

