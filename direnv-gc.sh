#!/usr/bin/env bash

set -eo pipefail
shopt -s nullglob

direnv_config=${DIRENV_CONFIG:-$(direnv status | grep DIRENV_CONFIG | cut -d" " -f2-)}
direnv_store=${DIRENV_STORE:-$direnv_config/store}

DAYS=10
DRY_RUN=0
VERBOSE=0

# parse arguments
for i in "$@"
do
    case $i in
	-d=*|--days=*)
	    DAYS="${i#*=}"
	    shift
	    ;;
	--dry)
	    DRY_RUN=1
	    shift
	    ;;
	-v|--verbose)
	    VERBOSE=1
	    shift
	    ;;
    esac
done

echo "$0: Removing any environments not loaded the last $DAYS days $([ $DRY_RUN -gt 0 ] && echo '[DRY RUN]')" >&2

# check if number
re='^[0-9]+$'
if ! [[ $DAYS =~ $re ]] ; then
    echo "$0: error: days not a number" >&2; exit 1
fi

if [ "$(uname)" == "Darwin" ]; then
    echo "Mac!"
    exit 1
else
    get_modified_timestamp() {
	echo $(stat -c %Y $1)
    }
    ts=$(date -d "-$DAYS days" +"%s")
fi

removed=0
bytes_removed=0
for FILE in $direnv_store/*; do
    [ $VERBOSE -gt 0 ] && echo "$0: Checking $FILE" >&2
    if [ "$(get_modified_timestamp $FILE)" -lt "$ts" ]; then
	if [[ -d "$(readlink $FILE)" ]]; then
	    echo "$0: $(readlink $FILE) - last modified: $(date -r $FILE)"
	    # gather metrics
	    read size _ < <(du -sb $(readlink $FILE))
	    ((bytes_removed=bytes_removed+size))
	    # remove
	    [ $DRY_RUN == 0 ] && rm -rf $(readlink $FILE)
	    [ $DRY_RUN == 0 ] && rm $FILE
	    ((removed=removed+1))
	else
	    echo "$0: Removing dead link $FILE -> $(readlink $FILE) "
	    [ $DRY_RUN == 0] rm $FILE
	fi
    fi
done

[ $bytes_removed -gt 0 ] && removed_in_mb=$(expr $bytes_removed / 1024 / 1024)
if [ $DRY_RUN -gt 0 ]; then
    echo "$0: Would clean up $removed environments for a total of $(echo $removed_in_mb)mb." >&2
else
    echo "$0: Cleaned up $removed environments for a total of $(echo $removed_in_mb)mb." >&2
fi

