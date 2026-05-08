#!/usr/bin/env bash
# Verify that every Go source file ships with the Apache 2.0 boilerplate.
set -euo pipefail

missing=0
while IFS= read -r f; do
  if ! head -5 "$f" | grep -q "Licensed under the Apache License"; then
    echo "missing license header: $f"
    missing=1
  fi
done < <(git ls-files '*.go' | grep -v zz_generated | grep -v _test.go)

if [ $missing -ne 0 ]; then
  exit 1
fi
echo "all source files carry the license header"
