#!/usr/bin/bash

set -euo pipefail

# source the utils file, get functions
source "./utils.sh"

PROJECT_DIR=$(get_project_root)
env_path=$(get_env "${PROJECT_DIR}")

if [[ -f ${env_path} ]]; then
    set -o allexport
    source "${env_path}"
    set +o allexport
else
    echo "env_path could not be evaluated, exiting"
    exit 1
fi

get_db_string() {
    local num_args="${#}"
    if [[ ${num_args} -lt 1 ]]; then
        echo "function get_db_string must be called with either \"dev\" or \"prod\""
        exit 1
    else
        if [[ "${1}" != "dev" ]] &&  [[ "${1}" != "prod" ]]; then
            echo "function get_db_string must be called with either \"dev\" or \"prod\""
            exit 1
        fi
        if [[ "${1}" == "dev" ]]; then
            echo "${DBSTRING}"
        fi
        if [[ "${1}" == "prod" ]]; then
            db_string="postgresql://${PROD_POSTGRES_USER}:${PROD_POSTGRES_PASSWORD}@localhost:5432/${PROD_DBNAME}?sslmode=${PROD_SSLMODE}&sslrootcert=${PROJECT_DIR}/global-bundle.pem"
            echo ${db_string}
        fi
    fi
}
