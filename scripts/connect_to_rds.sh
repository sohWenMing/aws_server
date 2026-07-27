#!/usr/bin/bash
set -euo pipefail

if [[ -f "./tunnel_to_rds.sh" ]]; then
    echo "tunnel script exists - continuing on to execution"
else
    echo "tunnel script does not exist, exiting script"
    exit 1
fi

# Sets up a tunnel to the ec2 instance that CAN connect to the postgres instance
# this will run the tunnel in detached mode - check instructions on getting rid of process
./tunnel_to_rds.sh

UTILS_FILE_PATH="./utils.sh"


# check that the utils file does exist, if so then load out all the utils functions
if [[ ! -f "${UTILS_FILE_PATH}" ]]; then
    echo "utils file does not exist - exiting"
    exit 1
else
   set -o allexport
   source "${UTILS_FILE_PATH}"
   set +o allexport
fi

PROJECT_ROOT=$(get_project_root)
ENV_PATH=$(get_env "${PROJECT_ROOT}")
if [[ ! -f "${ENV_PATH}" ]]; then
    echo ".env file does not exist - exiting"
    exit 1
else
    set -o allexport
    source "${ENV_PATH}"
    set +o allexport
fi

HOST="localhost"
PORT="5432"
SSLMODE="verify-ca"
SSLROOTCERTPATH="${PROJECT_ROOT}/global-bundle.pem"

# connect using psql - not that PGPASSWORD should be set as an environment variable
# so that it will not turn up in the logs
PGPASSWORD="${PROD_POSTGRES_PASSWORD}" psql \
    "host=${HOST} \
    port=${PORT} \
    dbname=${PROD_DBNAME} \
    user=${PROD_POSTGRES_USER} \
    sslmode=${SSLMODE} \
    sslrootcert=${SSLROOTCERTPATH}"
