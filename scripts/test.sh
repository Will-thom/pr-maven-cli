#!/usr/bin/env sh
set -eu

: "${GOCACHE:=$(pwd)/.gocache}"
: "${GOMODCACHE:=$(pwd)/.gomodcache}"
export GOCACHE GOMODCACHE

go test ./...

if [ "${PRMAVEN_RACE:-0}" = "1" ]; then
  go test -race ./...
fi

if [ "${PRMAVEN_COVERAGE:-0}" = "1" ]; then
  go test -coverprofile=coverage.out ./...
  coverage_output="$(go tool cover -func=coverage.out)"
  printf '%s\n' "$coverage_output"
  coverage="$(printf '%s\n' "$coverage_output" | awk '/^total:/ { gsub("%", "", $3); print $3 }')"
  minimum="${PRMAVEN_MIN_COVERAGE:-70}"
  awk -v coverage="$coverage" -v minimum="$minimum" 'BEGIN { if ((coverage + 0) < (minimum + 0)) exit 1 }'
fi
