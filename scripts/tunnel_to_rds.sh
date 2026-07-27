#!/usr/bin/bash
set -euo pipefail

source "./utils.sh"
PROJECT_DIR=$(get_project_root)
env_path=$(get_env "${PROJECT_DIR}")

# exports all environment variables that are sourced
set -o allexport
source "${env_path}"
set +o allexport
RDS_HOST_AND_PORT="${PROD_DB_HOST}:${PROD_DB_PORT}"
PORTLINK="5432:${RDS_HOST_AND_PORT}"
ssh -f -N -o ExitOnForwardFailure=yes -J jumphost -L "${PORTLINK}" private-a < /dev/null

sleep 1

TUNNELPID=$(pgrep -f "${PORTLINK}")
echo "${TUNNELPID}"
