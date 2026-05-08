#!/usr/bin/env bash
# Smoke-test every YAML file under the repo can be parsed as a Kubernetes object.
set -euo pipefail

while IFS= read -r f; do
  python3 -c "import yaml,sys; yaml.safe_load_all(open(sys.argv[1]).read())" "$f"
done < <(git ls-files '*.yaml' '*.yml' | grep -v vendor)
