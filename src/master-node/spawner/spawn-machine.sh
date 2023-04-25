#!/bin/bash

read -r cluster

if [ "$cluster" = "" ]
then
    cluster=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
fi

name="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

echo $name

container_id=$(docker run --expose 3000 -td --name $name alpine:3.14)

echo "machine $container_id running on port 3000\n"

existing_cluster=$(sqlite3 db "select name from cluster where name =\"$cluster\"")

echo $existing_cluster

cl=""

if [ "$existing_cluster" = "" ]
then
    cl=$cluster
    echo "$cl" | sudo sh master-node/etcd/cluster/create-cluster.sh
else
    cl=$existing_cluster
fi


echo "$cl:$container_id" | sudo sh master-node/etcd/create-machine.sh