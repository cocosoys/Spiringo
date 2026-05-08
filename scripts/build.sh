#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
OUT_DIR="${OUT_DIR:-$ROOT_DIR/bin}"
GOOS_VALUE="${GOOS:-}"
GOARCH_VALUE="${GOARCH:-amd64}"
CGO_ENABLED_VALUE="${CGO_ENABLED:-1}"

mkdir -p "$OUT_DIR"

if [ -n "$GOOS_VALUE" ]; then
  export GOOS="$GOOS_VALUE"
fi
export GOARCH="$GOARCH_VALUE"
export CGO_ENABLED="$CGO_ENABLED_VALUE"

go build -trimpath -o "$OUT_DIR/spiringo" "$ROOT_DIR/cmd/spiringo"
go build -trimpath -o "$OUT_DIR/spiringo-cli" "$ROOT_DIR/cmd/spiringo-cli"
go build -trimpath -o "$OUT_DIR/spiringo-serverless" "$ROOT_DIR/cmd/spiringo-serverless"
