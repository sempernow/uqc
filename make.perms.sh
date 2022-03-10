#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile : make perms 
# -----------------------------------------------------------------------------

perms(){
	echo "=== chmod 0600 @ all FILEs under '$1'"
	find "$1" -type f -execdir /bin/bash -c 'chmod 0600 "$@"' _ {} \+
}

perms assets

echo "=== chmod 0660 @ all *.{md,html,png,...} FILEs under PWD"
find . -type f -iname '*.md' -execdir /bin/bash -c 'chmod 0660 "$@"' _ {} \+ &
find . -type f -iname '*.html' -execdir /bin/bash -c 'chmod 0660 "$@"' _ {} \+ &
find . -type f -iname '*.json' -execdir /bin/bash -c 'chmod 0660 "$@"' _ {} \+ &
find . -type f -iname '*.log' -execdir /bin/bash -c 'chmod 0660 "$@"' _ {} \+ &
echo "=== chmod 0444 @ all 'LICENSE' files under PWD"
find . -type f -iname 'LICENSE' -execdir /bin/bash -c 'chmod 0444 "$@"' _ {} \+ &

echo "=== chmod 0774 @ all *.sh FILEs under PWD"
find . -type f -iname '*.sh' -execdir /bin/bash -c 'chmod 0774 "$@"' _ {} \+
echo "=== chmod 0774 @ all *.go FILEs under PWD"
find . -type f -iname '*.go' -execdir /bin/bash -c 'chmod 0774 "$@"' _ {} \+

echo '=== chmod 0755 @ all DIRs under PWD'
find . -type d -execdir /bin/bash -c 'chmod 0755 "$@"' _ {} \+ &

sleep 2
# Wait until all (background) proceses (@ '-execdir') complete
while [[ $(ps aux |grep -- -execdir |grep -v grep |awk 'NR == 1 {print $2}') ]]; do sleep 2; done


exit 0

