#!/bin/bash

if [ ! -f "./master-node/watcher/events.txt" ]
then
    touch "./master-node/watcher/events.txt"
fi

docker events --filter "type=container" --format "{{.ID}} {{.Type}} {{.Status}}" > ./master-node/watcher/events.txt &

tail -f ./master-node/watcher/events.txt | xargs -r -L1 sh ./master-node/watcher/handle-events.sh