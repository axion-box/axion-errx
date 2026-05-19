#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"

packages=()
while IFS= read -r line; do
  [[ -n "${line}" ]] || continue
  packages+=("${line}")
done < <(go list ./... || true)

if [[ ${#packages[@]} -eq 0 ]]; then
  echo "[run_unit_tests] 跳过: 当前仓库没有可执行 Go 包"
  exit 0
fi

echo "[run_unit_tests] running unit packages"
go test "${packages[@]}"
