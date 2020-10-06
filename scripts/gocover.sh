#!/usr/bin/env bash
# Usage: scripts/gocover.sh
#
# Coverage test

set -e

echo "mode: count" > coverage.txt

go test -v -covermode=count -coverprofile=profile.txt ./services/user/...
grep -v "mode: count" >> coverage.txt profile.txt

go test -v -covermode=count -coverprofile=profile.txt ./services/group/...
grep -v "mode: count" >> coverage.txt profile.txt

bash <(curl -s https://codecov.io/bash)

rm -rf ./coverage.txt
rm -rf ./profile.txt