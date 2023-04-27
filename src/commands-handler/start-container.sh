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

echo "starting machine..."
sudo bash ./master-node/spawner/spawn-machine.sh $1