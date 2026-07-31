#!/usr/bin/bash
set -euo pipefail

source "./utils.sh"
PROJECT_DIR=$(get_project_root)
BINARY_PATH="${PROJECT_DIR}/aws_server"
ENV="PROD" ${BINARY_PATH}
