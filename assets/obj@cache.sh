#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  List cache files per object 
#
#  ARG: post(default) | tags | categories | users
# -----------------------------------------------------------------------------

printf "\n%s  %s\n\n" "=== cache @ '*_${1:-post}.json'" '($1) {post(default), tags, categories, users}'
find ./../cache -iname "*_${1:-post}.json" -printf "%p\n" # |xargs -IX rm X

exit 

