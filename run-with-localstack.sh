#!/bin/bash
set -euo pipefail

BASE_DIR=$(cd $(dirname ${0}); pwd)

echo $BASE_DIR


if [ ! -f "${BASE_DIR}/.env" ]; then
    echo "[ERROR] .env file not found in ${BASE_DIR}" >&2
    exit 1
fi

export $(cat "${BASE_DIR}/.env" | sed 's/#.*//g' | xargs)

go run cmd/main.go