#!/bin/bash

docker_id=$1

status=$(docker inspect -f '{{.State.Status}}' $docker_id)
if [ "$status" = "" ]
then
    echo "no such container in this node" >&2
    exit -1
else
    if [ "$status" = "false" ]
    then
        echo "container $docker_id went down"
        exit -1
    fi

    echo "$status"
fi
