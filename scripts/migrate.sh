#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
ENV_NAME="${APP_ENV:-${ENV:-local}}"
CONFIG_DIR="${CONFIG_DIR:-$ROOT_DIR/configs}"

cd "$ROOT_DIR"
go run ./cmd/spiringo migrate up -env "$ENV_NAME" -config "$CONFIG_DIR"
