#!/bin/bash

id=$1

status=$(docker inspect -f '{{.State.Status}}' $id)
if [ "$status" = "" ]
then
    echo "no such container in this node" >&2
    exit -1
else
    if [ "$status" = "false" ]
    then
        echo "container $id went down"
        exit -1
    fi

    echo "$status"
fi
