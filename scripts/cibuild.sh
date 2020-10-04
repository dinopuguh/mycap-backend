#!/usr/bin/env bash
# Usage: scripts/cibuild
#
# Run tests

set -e
scripts/goget.sh
scripts/gomigrate.sh
scripts/gocover.sh