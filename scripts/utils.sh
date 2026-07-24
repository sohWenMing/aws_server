#!/usr/bin/bash

set -euo pipefail
get_project_root() {
    SCRIPT_DIR=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)
    # echo "SCRIPT_DIR: ${SCRIPT_DIR}"
    PROJECT_DIR="--"
    dir=${SCRIPT_DIR}

    while [[ "${dir}" != "/" ]]; do
        curEval="${dir}/go.mod"
        if [[ -f "${dir}/go.mod" ]]; then
            PROJECT_DIR=$(cd $(dirname ${curEval}) && pwd)
            break
        else
            dir=$(dirname "${dir}")
        fi
    done

    if [[ "${PROJECT_DIR}" == "--" ]]; then
        echo "project root could not be found, exiting with error"
        exit 1
    fi
    echo ${PROJECT_DIR}
}
