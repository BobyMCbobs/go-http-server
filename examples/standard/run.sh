#!/bin/bash

FOLDER="$(dirname "$(realpath "$0")")"
cd "$(git rev-parse --show-toplevel)"

export APP_SERVE_FOLDER="$FOLDER/folder"
go run .
