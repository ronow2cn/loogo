#!/bin/bash

export GOPATH=$(pwd)

# -----------------------------------------------
if [ $# -lt 1 ]; then
    echo installing gate ...
    go install gate

    echo installing game ...
    go install game

    echo installing client...
    go install client

    echo installing switcher...
    go install switcher

    echo installing auth...
    go install auth
else
    case $1 in
        gate)
            echo installing gate ...
            go install gate
            ;;
        client)
            echo installing client ...
            go install client
            ;;
        game)
            echo installing game ...
            go install game
            ;;
        switcher)
            echo installing switcher ...
            go install switcher
            ;;
        auth)
            echo installing auth ...
            go install auth
            ;;
        *)
            echo "Usage: build.sh gate|client or ...[build all]"
    esac
fi

echo copy config.json to bin...
cp config.json bin/

echo -e "\033[32mDone\033[0m"
