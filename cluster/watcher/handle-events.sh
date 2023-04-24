#!/bin/bash

id=$1
type=$2
status=$3

echo "id=$id;type=$type;status=$status"

if [ "$id" = "" ]
then
    exit 1
fi


if [ "$status" = "health_status" ]
then
    sh ./cluster/watcher/restart-container.sh $id
    exit 1
fi

if [ "$status" = "die" ]
then
    sh ./cluster/watcher/restart-container.sh $id
    exit 1
fi

if [ "$status" = "oom" ]
then
    sh ./cluster/watcher/restart-container.sh $id
    exit 1
fi
if [ "$status" = "start" ]
then
    #sh ./cluster/watcher/restart-container.sh $id
    exit 1
fi