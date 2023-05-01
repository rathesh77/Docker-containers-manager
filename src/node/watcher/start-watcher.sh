#!/bin/bash

if [ ! -f "./watcher/events.txt" ]
then
    touch "./watcher/events.txt"
fi

docker events --filter "type=container" --format "{{.ID}} {{.Type}} {{.Status}}" > ./watcher/events.txt &

tail -f ./watcher/events.txt | xargs -r -L1 sh ./watcher/handle-events.sh