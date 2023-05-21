#!/bin/bash

id="$1-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

image=$2

#$args="$3"
#container_id=$(docker run -td --name $name alpine:3.14)

pod_network="$1-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

docker network create $pod_network


container_id=$(docker run --network "$pod_network" -td --name "$id" $3 -p 8082:8080 "$image")


status="$(docker container inspect -f '{{.State.Running}}' $id)"

if [ "$status" = "false" ]
then
    docker network rm $pod_network
    echo "deleting container $id"
    docker rm $id
    echo "container failed to start" >&2
    exit -1
fi

if [ "$container_id" = "" ]
then
    docker network rm $pod_network
    echo "deleting container $id"
    docker rm $id
    echo "container failed to start" >&2
    exit -1
fi

echo "$container_id:$pod_network"