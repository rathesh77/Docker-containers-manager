#!/bin/bash

read -r namespace

if [ $namespace = "" ]
then
    $namespace=$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')
fi

name="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

echo $name

container_id=$(docker run --expose 3000 -td --name $name alpine:3.14)

echo "machine $container_id running on port 3000\n"

# jointure MANQUANT entre le conteneur et le cluster

echo "$namespace" | sudo sh cluster/etcd/namespace/create-namespace.sh
echo "$container_id" | sudo sh cluster/etcd/set-machine.sh