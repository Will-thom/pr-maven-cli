#!/usr/bin/env sh
set -eu

: "${GOCACHE:=$(pwd)/.gocache}"
: "${GOMODCACHE:=$(pwd)/.gomodcache}"
export GOCACHE GOMODCACHE

unformatted="$(gofmt -l ./cmd ./pkg)"
if [ -n "$unformatted" ]; then
  printf '%s\n' "$unformatted"
  printf '%s\n' "Go files need formatting. Run gofmt before committing." >&2
  exit 1
fi

go vet ./...
go test ./...
