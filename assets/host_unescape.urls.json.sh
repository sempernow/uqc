#!/usr/bin/env bash

[[ -f $1 ]] && {
    sed -i 's#\\/#/#g' "$1"
    true
} || {
    find ./../cache -type f -iname '*.json' -exec sed -i 's#\\/#/#g' {} \;
}