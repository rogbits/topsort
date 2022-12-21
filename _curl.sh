#!/usr/bin/env bash
set -e
port=8080

for file in fixtures/*; do
  echo "$file" "--------------------"
  echo "INPUT:"
  cat "$file"
  echo
  echo "OUTPUT"
  curl -s "localhost:$port/sort" -d @"$file"
  echo
  echo
done
