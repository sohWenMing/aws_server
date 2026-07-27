#!/usr/bin/bash
arg=$1
echo "arg: ${arg}"
source "./build_db_string.sh"

db_string=$(get_db_string ${arg})
echo "db_string: ${db_string}"
