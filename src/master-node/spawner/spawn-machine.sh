#!/bin/bash

read -r cluster

if [ "$cluster" = "" ]
then
    cluster=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
fi

name="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

existing_cluster=$(sqlite3 db "select name from cluster where name =\"$cluster\"")

cl=""

if [ "$existing_cluster" = "" ]
then
    cl=$cluster
    echo "$cl" | sudo sh master-node/etcd/cluster/create-cluster.sh
else
    cl=$existing_cluster
fi

network=$(sqlite3 db "select network from cluster where name =\"$cluster\"")

#container_id=$(docker run --network="$network" --expose 3000 -td --name $name alpine:3.14)

container_id=$(docker run -td --network="$network" --log-driver json-file --log-opt max-file=10 --log-opt max-size=2m -v /etc/localtime:/etc/localtime:ro -v ./config:/root/.stash -v ./data:/data -v ./metadata:/metadata -v ./cache:/cache -v ./blobs:/blobs -v ./generated:/generated --name $name -p 9999:9999 -e STASH_STASH=/data/ -e STASH_GENERATED=/generated/ -e STASH_METADATA=/metadata/ -e STASH_CACHE=/cache/ -e STASH_PORT=9999 stashapp/stash:latest)

echo "machine $container_id running on port 3000\n"

echo "$cl:$container_id" | sudo sh master-node/etcd/create-machine.sh