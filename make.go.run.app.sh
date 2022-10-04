#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile recipes for : go run ...  
# -----------------------------------------------------------------------------

export APP_CLIENT_PASS="$(cat ./assets/.env/app.env \
    |grep APP_CLIENT_PASS |awk -F '=' '{print $2}'\
)"


token() {
    go run ./app/cli token
}

uptkn() { # UpsertMsgByTkn
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    [[ $1 ]] && json="$1" || json="$(<assets/msg.json)"
    [[ $2 ]] && mid="$2"  || mid="$(uuid -v 5 ns:OID  /2022/09/15/uqrate-client-test)"
    [[ $3 ]] && tkn="$3"  || tkn="$(go run ./app/cli token |jq -Mr .body)"

    #go run ./app/cli uptkn "$json" "$mid" "${tkn}" "${APP_CHANNEL_SLUG}"
    go run ./app/cli uptkn "$json" "${tkn}" "${APP_CHANNEL_SLUG}"
}

upkey() { # UpsertMsgByKey
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    [[ $1 ]] && json="$1" || json="$(<assets/msg.json)"
    [[ $3 ]] && key="$2"  || key="$(cat ./assets/keys/uqrate.${APP_CHANNEL_SLUG}.json |jq -Mr .key)"

    echo "key : name: '${key%.*}'"
    go run ./app/cli upkey "$json" "${key}"
}

wpfetch() { # Fetch from a WP posts endpoint and dump JSON response to file
    url=$1
    fname=${url#*//}
    fname=${fname%%/*}
    obj=${url##*/}
    obj=${obj%\?*}
    go run ./app/cli wpfetch $url > assets/wp/${fname}.${obj}.json
}

wpuptkn() { 
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }

    url=$1
    fname=${url#*//}
    fname=${fname%%/*}
    [[ $2 ]] && tkn="$2"  || tkn="$(go run ./app/cli token |jq -Mr .body)"
    [[ $2 ]] && slug="$3" || slug="${APP_CHANNEL_SLUG}"
    [[ $3 ]] && os="$4"   || os="${APP_CLIENT_USER}${APP_CHANNEL_SLUG}"

    go run ./app/cli wpupkey "assets/wp/posts.${fname}.json" "$tkn" "${slug}" "$os"
}
wpupkey() { 
    [[ $(type -t jq) ]] || { echo 'REQUIREs jq utility'; exit 0; }
    export key="$(cat ./assets/keys/uqrate.${APP_CHANNEL_SLUG}.json |jq -Mr .key)"
    export chn_id="$(cat ./assets/keys/uqrate.${APP_CHANNEL_SLUG}.json |jq -Mr .chn_id)"

    url=$1
    fname=${url#*//}
    fname=${fname%%/*}
    go run ./app/cli wpupkey "assets/wp/posts.${fname}.json" "$key" "$chn_id"
}

cli() { # Any
    go run ./app/cli "$@"
}

"$@"

