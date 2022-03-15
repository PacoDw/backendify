#!/usr/bin/env bash

COMMAND_TOOLS="golangci-lint"
MISSING_TOOLS=false

echo ""
echo "==> Checking if there is any missing dependencies..."

for cmd in ${COMMAND_TOOLS} ; do
    if ! command -v bin/${cmd} &> /dev/null ; then
        echo "  -> ${cmd} is required"
        MISSING_TOOLS=true
    fi
done

if ${MISSING_TOOLS} = true ; then
    exit 1
else
    echo "The tools are already installed!";
fi

exit 0
