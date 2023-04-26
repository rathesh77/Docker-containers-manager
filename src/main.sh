#!/bin/bash

sudo ./master-node/watcher/start-watcher.sh &

while true;
do
    read -e line

    command=$(echo "$line" | sed "s/ .*//")

    args=$(echo "$line" | sed -r "s/([a-z\-]+ ){1}//")

    if [ "$args" != "$command" ]
    then
        case $command in

            start-container)
                ./commands-handler/start-container.sh $args
            ;;

        esac
    fi
done;