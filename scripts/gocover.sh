#!/usr/bin/env bash
# Usage: scripts/gocover.sh
#
# Coverage test

set -e

echo "mode: set" > coverage.out

go test -v -covermode=set -coverprofile=profile.out ./services/user/...
grep -v "mode: set" >> coverage.out profile.out

go test -v -covermode=set -coverprofile=profile.out ./services/group/...
grep -v "mode: set" >> coverage.out profile.out

echo "$GOPATH/bin/goveralls -coverprofile=coverage.out -service=travis-pro -repotoken $COVERALLS_TOKEN"

rm -rf ./coverage.out
rm -rf ./profile.out