#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile recipes for : go run ...  
# -----------------------------------------------------------------------------

export APP_ASSETS=${PATH_HOST_ROOT}/assets
export APP_CACHE=${PATH_HOST_ROOT}/cache

export APP_CLIENT_PASS="$(cat ${APP_ASSETS}/.env/app.${APP_SERVICE_HOST}.env \
    |grep APP_CLIENT_PASS |awk -F '=' '{print $2}'\
)"
export APP_CLIENT_KEY="$(cat ${APP_ASSETS}/.env/app.${APP_SERVICE_HOST}.env \
    |grep APP_CLIENT_KEY |awk -F '=' '{print $2}'\
)"
export APP_SITES_PASS="$(cat ${APP_ASSETS}/.env/app.${APP_SERVICE_HOST}.env \
    |grep APP_SITES_PASS |awk -F '=' '{print $2}'\
)"

# echo "key : $APP_CLIENT_KEY"
siteslist() {
    go run ./app/cli siteslist
}
updateusers() {
    go run ./app/cli updateusers
}
upsertchns() {
    go run ./app/cli upsertchns
}
upsertposts() {
    go run ./app/cli upsertposts
}
token() {
    go run ./app/cli token
}
key() {
    go run ./app/cli key "$APP_CHANNEL_ID"
}
purgecachetkns() {
    go run ./app/cli purgecachetkns 
    #find "${APP_CACHE:-cache}" -iname 'tkn.*' -exec rm {} \;
}
purgecacheposts() {
    go run ./app/cli purgecacheposts
    #find "${APP_CACHE:-cache}" -iname '*_posts.json' -exec rm {} \;
}

uptkn() { # UpsertMsgByTkn
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    [[ $1 ]] && json="$1" || json="$(<${APP_ASSETS}/msg.json)"
    [[ $2 ]] && mid="$2"  || mid="$(uuid -v 5 ns:OID  /2022/09/15/uqrate-client-test)"
    [[ $3 ]] && tkn="$3"  || tkn="$(go run ./app/cli token)"

    #go run ./app/cli uptkn "$json" "$mid" "${tkn}" "${APP_CHANNEL_SLUG}"
    go run ./app/cli uptkn "$json" "${tkn}" "${APP_CHANNEL_SLUG}"
}

upkey() { # UpsertMsgByKey
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    [[ $1 ]] && json="$1" || json="$(<${APP_ASSETS}/msg.json)"
    [[ $3 ]] && key="$2"  || key="$(cat ${APP_CACHE}/keys/key.${APP_CHANNEL_ID}.json |jq -Mr .key)"

    # echo "key : name: '${key%.*}'"
    go run ./app/cli upkey "$json" "${key}"
}

wpfetch() { # Fetch from a WP posts endpoint and dump JSON response to file
    url=$1              # https://wp.site/wp-json/wp/v2/posts?author=7
    fname=${url#*//}    # wp.site/wp-json/wp/v2/posts?author=7
    fname=${fname%%/*}  # wp.site
    obj=${url##*/}      # posts?author=7
    obj=${obj%\?*}      # posts
    go run ./app/cli wpfetch $url > ${APP_ASSETS}/wp/${fname}_${obj}.json
}

wpuptkn() { 
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }

    url=$1
    fname=${url#*//}
    fname=${fname%%/*}
    [[ $2 ]] && tkn="$2"  || tkn="$(go run ./app/cli token |jq -Mr .body)"
    [[ $2 ]] && slug="$3" || slug="${APP_CHANNEL_SLUG}"
    [[ $3 ]] && os="$4"   || os="${APP_CLIENT_USER}${APP_CHANNEL_SLUG}"

    go run ./app/cli wpupkey "${APP_ASSETS}/wp/posts.${fname}.json" "$tkn" "${slug}" "$os"
}
wpupkey() { 
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    export key="$(cat ${APP_ASSETS}/keys/uqrate.${APP_CHANNEL_SLUG}.json |jq -Mr .key)"
    export chn_id="$(cat ${APP_ASSETS}/keys/uqrate.${APP_CHANNEL_SLUG}.json |jq -Mr .chn_id)"

    url=$1
    fname=${url#*//}
    fname=${fname%%/*}
    go run ./app/cli wpupkey "${APP_ASSETS}/wp/posts.${fname}.json" "$key" "$chn_id"
}

cli() { # Any
    go run ./app/cli "$@"
}

"$@"

