#!/usr/bin/env bash
set -euo pipefail

ARCUBASE_ROOT="${1:-${ARCUBASE_ROOT:-}}"
if [[ -z "${ARCUBASE_ROOT}" ]]; then
  echo "ARCUBASE_ROOT is required" >&2
  exit 1
fi

UPSTREAM_FILE="${ARCUBASE_ROOT}/platform/types/controlplane/controlplane_swagger.json"
LOCAL_FILE="$(cd "$(dirname "$0")/.." && pwd)/contracts/controlplane_swagger.json"

if [[ ! -f "${UPSTREAM_FILE}" ]]; then
  echo "Upstream contract not found: ${UPSTREAM_FILE}" >&2
  exit 1
fi
if [[ ! -f "${LOCAL_FILE}" ]]; then
  echo "Local contract snapshot not found: ${LOCAL_FILE}" >&2
  exit 1
fi

if ! diff -u "${LOCAL_FILE}" "${UPSTREAM_FILE}"; then
  echo "Control-plane contract drift detected" >&2
  exit 1
fi

echo "Control-plane contract is in sync"
