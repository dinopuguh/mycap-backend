#!/usr/bin/env bash
# Usage: scripts/cibuild.sh
#
# Run tests

set -e
scripts/goget.sh
scripts/gomigrate.sh
scripts/goseed.sh
scripts/gocover.sh