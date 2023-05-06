#!/bin/bash

id="$1-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

#container_id=$(docker run -td --name $name alpine:3.14)

name=$1

pod_network="net-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

docker network create $pod_network

container_id=$(docker run \
    --network $pod_network \
    -td \
    --log-driver \
    json-file \
    --log-opt max-file=10 \
    --log-opt max-size=2m \
    -v /etc/localtime:/etc/localtime:ro \
    -v ./config:/root/.stash \
    -v ./data:/data \
    -v ./metadata:/metadata \
    -v ./cache:/cache \
    -v ./blobs:/blobs \
    -v ./generated:/generated \
    --name $id \
    -p 9999:9999 \
    -e STASH_STASH=/data/ \
    -e STASH_GENERATED=/generated/ \
    -e STASH_METADATA=/metadata/ \
    -e STASH_CACHE=/cache/ \
    -e STASH_PORT=9999 \
    stashapp/stash:latest)


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

echo "$id:$pod_network"