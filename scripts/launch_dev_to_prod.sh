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

source "./utils.sh"
PROJECT_DIR=$(get_project_root)
MAIN_PATH="${PROJECT_DIR}/main.go"
ENV="DEV_TO_PROD" go run "${MAIN_PATH}"
