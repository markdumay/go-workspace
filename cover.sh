#!/bin/sh

set -e

arg="$1"
[ -z "$arg" ] && arg="..."
arg="./${arg}"

echo "Generating coverage report"
go test "${arg}" -coverprofile=cover.out || exit
echo

echo "Running coverage tool"
go tool cover -html=cover.out