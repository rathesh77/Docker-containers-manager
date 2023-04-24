#!/bin/bash

read -r namespace

if [ "$namespace" = "" ]
then
    namespace=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
fi

name="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

echo $name

container_id=$(docker run --expose 3000 -td --name $name alpine:3.14)

echo "machine $container_id running on port 3000\n"

existing_namespace=$(sqlite3 db "select name from namespace where name =\"$namespace\"")

echo $existing_namespace

ns=""

if [ "$existing_namespace" = "" ]
then
    ns=$namespace
    echo "$ns" | sudo sh cluster/etcd/namespace/create-namespace.sh
else
    ns=$existing_namespace
fi


echo "$ns:$container_id" | sudo sh cluster/etcd/set-machine.sh