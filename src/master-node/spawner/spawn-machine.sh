#!/bin/bash

id="$1-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

node_id=$(sqlite3 db "select id from node limit 1")

#container_id=$(docker run -td --name $name alpine:3.14)

name=$1

container_id=$(docker run \
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
    stashapp/stash:latest \
    && sudo sh master-node/etcd/machine/create-machine.sh $node_id $id $name)


status="$(docker container inspect -f '{{.State.Running}}' $id)"

if [ "$status" = "false" ]
then
    echo "deleting container $id"
    docker rm $id
fi
