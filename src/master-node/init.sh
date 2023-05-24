#!/bin/bash

cluster_id=$(sqlite3 db "select id from cluster limit 1")
if [ "$cluster_id" = "" ]
then
    sudo bash ./spawner/spawn-cluster.sh
    #sudo bash ./spawner/spawn-node.sh
else
    echo "cluster already created"
fi

sudo bash ./controllers/discover-network.sh

cd ./api-server
go run . &

cd ..
while :
do
    echo "WATCH !"
    sudo bash ./controllers/watch-events.sh
    echo "sleep for 10 seconds !"
    sleep 10
done


