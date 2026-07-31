#!/usr/bin/bash
set -euo pipefail

source "./utils.sh"
PROJECT_DIR=$(get_project_root)
MAIN_PATH="${PROJECT_DIR}/main.go"
ENV="PROD" go run "${MAIN_PATH}"
