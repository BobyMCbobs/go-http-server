#!/bin/bash

FOLDER="$(dirname "$(realpath "$0")")"
cd "$(git rev-parse --show-toplevel)"

export APP_SERVE_FOLDER="$FOLDER/folder" \
    APP_TEMPLATE_MAP_PATH="$FOLDER/template-map.yaml" \
    APP_VUEJS_HISTORY_MODE=true
go run .
