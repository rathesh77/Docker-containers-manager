#!/bin/bash

if [ "$1" = "" ]
then
    echo "cluster name required"
    exit -1
fi


if [ "$2" != "" ]
then
    echo "too many args, only one needed"
    exit -1
fi

echo "$1" | ./master-node/spawner/spawn-machine.sh