#!/usr/bin/bash
set -euo pipefail
echo "running command: ssh -J jumphost private-a"
# ssh -J jumphost private-a

source "./utils.sh"
PROJECT_DIR=$(get_project_root)
env_path="--"
if [[ -f "${PROJECT_DIR}/.env" ]]; then
    env_path="${PROJECT_DIR}/.env"
else
    echo ".env file could not be found. exiting"
    exit 1
fi

set -o allexport
source "${env_path}"
set +o allexport
SSL_ROOT_CERT_PATH="${PROJECT_DIR}/global-bundle.pem"
DBSTRING="postgresql://${PROD_POSTGRES_USER}:${PROD_POSTGRES_PASSWORD}@${PROD_DB_HOST}:${PROD_DB_PORT}/${PROD_DBNAME} sslmode=${PROD_SSLMODE} sslrootcert=${SSL_ROOT_CERT_PATH}"
echo "DBSTRING: ${DBSTRING}"
