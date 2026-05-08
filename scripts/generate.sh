#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
MODULE_NAME="${MODULE:-}"
MODEL_NAME="${MODEL:-}"
TYPE="${TYPE:-}"

cd "$ROOT_DIR"

case "$TYPE" in
  module)
    if [ -z "$MODULE_NAME" ]; then
      echo "MODULE is required for TYPE=module" >&2
      exit 2
    fi
    go run ./cmd/spiringo-cli module "$MODULE_NAME"
    ;;
  crud)
    if [ -z "$MODULE_NAME" ] || [ -z "$MODEL_NAME" ]; then
      echo "MODULE and MODEL are required for TYPE=crud" >&2
      exit 2
    fi
    go run ./cmd/spiringo-cli crud "$MODULE_NAME" "$MODEL_NAME"
    ;;
  payment-channel)
    if [ -z "$MODULE_NAME" ]; then
      echo "MODULE is required for TYPE=payment-channel" >&2
      exit 2
    fi
    go run ./cmd/spiringo-cli payment-channel "$MODULE_NAME"
    ;;
  oauth-provider)
    if [ -z "$MODULE_NAME" ]; then
      echo "MODULE is required for TYPE=oauth-provider" >&2
      exit 2
    fi
    go run ./cmd/spiringo-cli oauth-provider "$MODULE_NAME"
    ;;
  *)
    echo "Usage examples:" >&2
    echo "  TYPE=module MODULE=order scripts/generate.sh" >&2
    echo "  TYPE=crud MODULE=order MODEL=Product scripts/generate.sh" >&2
    echo "  TYPE=payment-channel MODULE=bank scripts/generate.sh" >&2
    echo "  TYPE=oauth-provider MODULE=github scripts/generate.sh" >&2
    exit 2
    ;;
esac
