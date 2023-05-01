#!/bin/bash

cluster_id=$(sqlite3 db "select id from cluster limit 1")
if [ "$cluster_id" = "" ]
then
    sudo bash ./spawner/spawn-cluster.sh
    sudo bash ./spawner/spawn-node.sh
else
    echo "cluster already created"
fi

#sudo sh ./master-node/watcher/start-watcher.sh &

cd ./api-server
go run .


