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

cp "${UPSTREAM_FILE}" "${LOCAL_FILE}"
echo "Synced contract to ${LOCAL_FILE}"
