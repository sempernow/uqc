#!/usr/bin/env bash

_sql() {
    # Requires root only; no trailing slash
    # https://www.foo.bar.com => foo.bar
    url="$@"
    obj=${url##*/}
    obj=${obj%\?*}
    obj=${obj/.com/}
    obj=${obj/www./}
    obj=${obj/.org/}
    obj=${obj/.net/}
    
    echo "=== @ obj  : '$obj'" 

    printf "%s\n" "('{UMEM,HOST}', 'proxy', pw_hash('proxy'), '${obj}', '${obj}', '${obj}@emx.unk')," >> $users

    printf "%s\n" "(
        (SELECT view_id FROM views WHERE vname = 'chn-view'), 
        (SELECT user_id FROM users WHERE handle = '${obj}'), 
        'Mirror', msgform_long(),'Mirror', 'Proxy', '$url'
    )," >> $channels

}
export -f _sql 

export users=host_insert.users.sql
export channels=host_insert.channels.sql

echo > $users
echo > $channels

[[ $1 ]] && {
    sql "$1"
    true
} || {
    
    printf "%s\n" 'INSERT INTO users (roles, about, pass_hash, handle, display, email) VALUES' >> $users
    printf "%s\n" 'INSERT INTO channels (view_id, owner_id, slug, msg_form, title, about, host_url) VALUES' >> $channels

    printf "%s\n" "$(cat host_urls)" |xargs -n 1 -I {} /bin/bash -c '_sql "$@"' _ {}

    printf "%s\n" 'ON CONFLICT DO NOTHING;' >> $users
    printf "%s\n" 'ON CONFLICT DO NOTHING;' >> $channels
}

exit
