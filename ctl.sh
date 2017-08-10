#!/bin/bash

SERVERS=(
)

_in_list()
{
    local val=$1
    local list=$2

    if [ -z "$list" ]; then
        return 0
    fi

    for v in $list; do
        if [ "$val" == "$v" ]; then
            return 0
        fi
    done
    return 1
}

_is_running()
{
    local cmd=$1
    local pid=$2

    if [ -d /proc/$pid ] && [[ "$(</proc/$pid/cmdline)" == *"${cmd// /}"* ]]; then
        return 0
    else
        return 1
    fi
}

_start()
{
    local names=$@

    # set limits
    ulimit -Sn unlimited

    # check running
    if ! _status no $names; then
        echo -e "\033[33mthere are still some servers running!\033[0m"
        _status yes $names
        exit 1
    fi

    # work
    for svr in "${SERVERS[@]}"; do

        svr_name=${svr%% *}
        svr_exec=${svr#* }

        if ! _in_list $svr_name "$names"; then
            continue
        fi

        echo -n "$(printf '%-32s' "starting $svr_name ...")"
        $svr_exec &>> $svr_name.err &
        pid=$!
        echo $pid > $svr_name.pid
        sleep 2
        if _is_running "$svr_exec" $pid; then
            echo -e "\033[32m[OK]\033[0m"
        else
            echo -e "\033[31m[Failed]\033[0m"
        fi
    done
}

_stop()
{
    local names=$@

    for ((i=${#SERVERS[@]}-1; i>=0; i--)); do

        svr=${SERVERS[i]}

        svr_name=${svr%% *}
        svr_exec=${svr#* }

        if ! _in_list $svr_name "$names"; then
            continue
        fi

        pidfile=$svr_name.pid
        if [ -f $pidfile ]; then
            echo -n "$(printf '%-32s' "stopping $svr_name ...")"

            #wait some time before stopping next
            sleep 1

            pid=$(<$pidfile)
            if _is_running "$svr_exec" $pid; then
                kill $pid
                while [ -d /proc/$pid ]; do
                    sleep 1
                done
            fi
            echo -e "\033[32m[STOPPED]\033[0m"
            rm $pidfile
        fi
    done
}

_status()
{
    local show=$1
    shift
    local names=$@
    local ret=0

    for svr in "${SERVERS[@]}"; do

        svr_name=${svr%% *}
        svr_exec=${svr#* }

        if ! _in_list $svr_name "$names"; then
            continue
        fi

        if [ -f $svr_name.pid ]; then
            pid=$(<$svr_name.pid)
            if _is_running "$svr_exec" $pid; then
                [ "$show" == "yes" ] && echo -e "$(printf '%-16s%-16s\033[32m[Running]\033[0m' $svr_name $pid)"
                ret=1
                continue
            fi
        fi
        [ "$show" == "yes" ] && echo -e "$(printf '%-32s\033[31m[Stopped]\033[0m' $svr_name)"
    done

    return $ret
}

# ---------------------------------------------------------

if [ -d bin ]; then
    wdir=bin
else
    wdir=.
fi

pushd $wdir > /dev/null

# import SERVERS
[ -f SERVERS ] && . SERVERS

cmd=$1
shift

case "$cmd" in
    start)
        _start $@
        ;;
    stop)
        _stop $@
        ;;
    restart)
        _stop $@
        _start $@
        ;;
    status)
        _status yes $@
        ;;
    *)
        echo "Usage: ctl.sh {start|stop|restart|status} [server...]"
        echo
        echo "    available servers:"
        for svr in "${SERVERS[@]}"; do
            svr_name=${svr%% *}
            echo "        $svr_name"
        done
esac

popd > /dev/null
