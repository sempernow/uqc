#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile recipes for : go run ...  
# -----------------------------------------------------------------------------

# All path settings at Makefile.settings are for the container environment,
# so all need resetting for this `go run` environment. Set them here:
# (Makefile isn't well suited for per-recipe variable settings.)

export APP_CLIENT_PASS="$(cat ./assets/.env/app.env \
    |grep APP_CLIENT_PASS |awk -F '=' '{print $2}'\
)"

token(){
    go run ./app/cli token
}

upsert(){
	export tkn="$(go run ./app/cli token |jq -Mr .body)"
    go run ./app/cli upsert "$(<assets/msg.json)" $(uuid -v 5 ns:OID  /2022/09/15/uqrate-client-test) \
        "${tkn}" "${APP_CHANNEL_SLUG}"
}

"$@"

