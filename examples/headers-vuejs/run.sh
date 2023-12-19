#!/bin/bash

FOLDER="$(dirname "$(realpath "$0")")"
cd "$(git rev-parse --show-toplevel)"

export APP_SERVE_FOLDER="$FOLDER/folder" \
    APP_HEADER_SET_ENABLE=true \
    APP_HEADER_MAP_PATH="$FOLDER/headers.yaml" \
    APP_TEMPLATE_MAP_PATH="$FOLDER/template-map.yaml" \
    APP_VUEJS_HISTORY_MODE=true
go run .
