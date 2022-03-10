#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  go run ...  any ./app/ from project root
# -----------------------------------------------------------------------------

# All path settings at Makefile.settings are for the container environment,
# so all need resetting for this `go run` environment. Set them here:
# (Makefile does not allow for per-recipe variable settings.)

export APP_CLIENT_PASS="$(cat ./assets/.env/app.env \
    |grep APP_CLIENT_PASS |awk -F '=' '{print $2}'\
)"


go run ./app/"$@"

