#!/usr/bin/env bash
set -e

if [ "$1" = "ci" ]; then
  echo using ci en
fi

mkdir -p coverage
go test -v -coverprofile coverage/cover.out ./...
if [ "$1" != "ci" ]; then
  go tool cover -html coverage/cover.out -o coverage/cover.html
  open coverage/cover.html
fi
