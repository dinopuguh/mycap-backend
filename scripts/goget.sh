#!/usr/bin/env bash
# Usage: script/goget.sh
#
# Downloads dependencies

set -e
go get ./...
go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls