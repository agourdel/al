#!/bin/bash
# Al go wrapper - allows changing directory in current shell
# Usage: algo <project_shortcut>

target_path=$(al go "$@" 2>/dev/null)
exit_code=$?

if [ $exit_code -eq 0 ] && [ -n "$target_path" ]; then
    cd "$target_path" || exit 1
else
    al go "$@"
    exit $?
fi
