#!/usr/bin/env bash
# Usage: scripts/cibuild.sh
#
# Run tests

set -e
scripts/gomigrate.sh
scripts/goseed.sh
scripts/gocover.sh