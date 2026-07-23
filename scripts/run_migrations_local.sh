#!/usr/bin/bash

set -euo pipefail

# get number of arguments passed when calling script
numArgs="$#"

# sentinel value - will be reset if operation can be found
operation="--"

if [[ ${numArgs} -lt 1 ]]; then
    echo "migration script must be called with at least one argument"
    exit 1
else
    operation="${1}"
fi

if [[ "${operation}" != "up" ]] && [[ "${operation}" != "down" ]]; then
    echo "operation for migration must be either \"up\" or \"down\""
    exit 1
fi

# gets the utility functions from ./utils.sh - sources get_project_root, function used to orientate to root
# of project
source "./utils.sh"

PROJECT_DIR=$(get_project_root)

env_path="${PROJECT_DIR}/.env"
if [[ -f ${env_path} ]]; then
    set -o allexport
    source "${env_path}"
    set +o allexport
fi

# set all the environment variables, before calling either the up or down migration depending on user input
# check documentation regarding env variables for goose - https://pressly.github.io/goose/documentation/environment-variables/
GOOSE_MIGRATION_DIR="${PROJECT_DIR}/${GOOSE_MIGRATION_DIR_SUFFIX}"
GOOSE_DBSTRING="${DBSTRING}"

echo "GOOSE_MIGRATION_DIR: ${GOOSE_MIGRATION_DIR}"
echo "GOOSE_DRIVER: ${GOOSE_DRIVER}"
echo "GOOSE_DBSTRING: ${GOOSE_DBSTRING}"

export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

goose "${operation}"
