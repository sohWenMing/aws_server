#!/usr/bin/bash

echo "launching application connecting to dev database..."
source "./utils.sh"

PROJECT_DIR=$(get_project_root)
MAIN_PATH="${PROJECT_DIR}/main.go"

ENV=DEV go run "${MAIN_PATH}"
