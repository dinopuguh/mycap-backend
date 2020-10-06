#!/usr/bin/env bash
# Usage: scripts/gomigrate.sh
#
# Run migrations

set -e
go run cmd/migrate/migrate.go