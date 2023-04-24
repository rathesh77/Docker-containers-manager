#!/bin/bash

if [ ! -f "./cluster/watcher/events.txt" ]
then
    touch "./cluster/watcher/events.txt"
fi

docker events --filter "type=container" --format "{{.ID}} {{.Type}} {{.Status}}" > ./cluster/watcher/events.txt &

tail -f ./cluster/watcher/events.txt | xargs -r -L1 sh ./cluster/watcher/handle-events.sh