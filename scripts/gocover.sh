#!/usr/bin/env bash
# Usage: scripts/gocover
#
# Coverage test

set -e

echo "mode: count" > coverage.out

go test -v -covermode=count -coverprofile=profile.out ./services/user/...
cat profile.out | grep -v "mode: count" >> coverage.out

go test -v -covermode=count -coverprofile=profile.out ./services/group/...
cat profile.out | grep -v "mode: count" >> coverage.out

$GOPATH/bin/goveralls -coverprofile=coverage.out -service=travis-pro -repotoken $COVERALLS_TOKEN

rm -rf ./coverage.out
rm -rf ./profile.out