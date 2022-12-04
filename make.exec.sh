#!/usr/bin/env bash
###############################################################################
# Makefile : make SVC : docker exec -it ...
# 
###############################################################################

export svc=$1
export ctnr=$(docker ps -q -f name=${svc} -f status=running | head -n 1)
docker exec -it $ctnr sh

exit 0
