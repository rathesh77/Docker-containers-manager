#!/bin/bash

id=$1

docker stop $id && sudo sh ./master-node/etcd/machine/delete-machine.sh $id