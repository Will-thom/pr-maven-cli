#!/usr/bin/env sh
set -eu

out_dir="${1:-dist}"
version="${2:-dev}"

: "${GOCACHE:=$(pwd)/.gocache}"
: "${GOMODCACHE:=$(pwd)/.gomodcache}"
export GOCACHE GOMODCACHE

mkdir -p "$out_dir"
go build -trimpath -ldflags="-s -w -X main.version=$version" -o "$out_dir/prmaven" ./cmd/prmaven
printf '%s\n' "$out_dir/prmaven"
