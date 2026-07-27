#!/usr/bin/bash

set -euo pipefail

# get number of arguments passed when calling script
numArgs="$#"

# sentinel value - will be reset if operation can be found
env="--"
operation="--"

if [[ ${numArgs} -lt 2 ]]; then
    echo "migration script must be called with at least two argument"
    exit 1
else
    env="${1}"
    operation="${2}"
fi

if [[ "${operation}" != "up" ]] && [[ "${operation}" != "down" ]]; then
    echo "usage: ./run_migrations_local.sh <dev/prod> <up/down>"
    echo "arg 1: operation for migration must be either \"up\" or \"down\""
    exit 1
fi
if [[ "${env}" != "dev" ]] && [[ "${env}" != "prod" ]]; then
    echo "usage: ./run_migrations_local.sh <dev/prod> <up/down>"
    echo "arg 2: env for migration must be either \"dev\" or \"prod\""
    exit 1
fi

# source build_db_string.sh - required to get utility function to get db_string
source "./utils.sh"
PROJECT_DIR=$(get_project_root)
BUILD_STRING_PATH="${PROJECT_DIR}/scripts/build_db_string.sh"
if [[ -f "${BUILD_STRING_PATH}" ]]; then
    source "${BUILD_STRING_PATH}"
else
    echo "build string path could not be found, exiting script"
    exit 1
fi


# set all the environment variables, before calling either the up or down migration depending on user input
# check documentation regarding env variables for goose - https://pressly.github.io/goose/documentation/environment-variables/
GOOSE_MIGRATION_DIR="${PROJECT_DIR}/${GOOSE_MIGRATION_DIR_SUFFIX}"

# set the dbstring - based on whether the environment passed in was dev or prod
GOOSE_DBSTRING=$(get_db_string ${env})

# Sets up a tunnel to the ec2 instance that CAN connect to the postgres instance
# this will run the tunnel in detached mode - check instructions on getting rid of process

TUNNELPID="--"
if [[ ${env} == "prod" ]]; then
    if [[ -f "${PROJECT_DIR}/scripts/tunnel_to_rds.sh" ]]; then
        echo "tunnel script exists - continuing on to execution"
        TUNNELPID=$("${PROJECT_DIR}"/scripts/tunnel_to_rds.sh)
    else
        echo "tunnel script does not exist, exiting script"
        exit 1
    fi
fi

echo "GOOSE_MIGRATION_DIR: ${GOOSE_MIGRATION_DIR}"
echo "GOOSE_DRIVER: ${GOOSE_DRIVER}"
echo "GOOSE_DBSTRING: ${GOOSE_DBSTRING}"

export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

goose "${operation}"

if [[ ${env} == "prod" ]] && [[ -n "${TUNNELPID}" ]] && [[ "${TUNNELPID}" != "--" ]]; then
    echo "killing tunnel process. PID: ${TUNNELPID}"
    kill "${TUNNELPID}"
fi
