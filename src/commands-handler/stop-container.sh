#!/bin/bash

if [ "$1" = "" ]
then
    echo "machine name required"
    exit -1
fi


if [ "$2" != "" ]
then
    echo "too many args, only one needed"
    exit -1
fi

sudo sh ./master-node/spawner/stop-machine.sh $1