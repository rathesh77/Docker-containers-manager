#!/bin/bash
echo 'ping' | nc -q  1  localhost 2000

response="$?"

if [ "$response" = "1" ];
then
    echo "init.sh not running"
    exit -1
fi

line=""

for word in "$@"; do line="$line $word"; done

line=$(echo "$line" | xargs)
echo "command: $line"

if [ "$line" = "" ]
then
    exit -1
fi

command=$(echo "$line" | sed "s/ .*//")

args=$(echo "$line" | sed -r "s/([a-z\-]+ ){1}//")

echo "command:$command"
echo "args:$args"
if [ "$args" != "$command" ]
then
    case $command in

        "start-container")
            sudo sh ./commands-handler/start-container.sh $args
        ;;

        "stop-container")
            sudo sh ./commands-handler/stop-container.sh $args
        ;;

    esac
fi