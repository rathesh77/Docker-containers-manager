#!/bin/bash

node_id=$1
container_docker_id=$2
pod_name=$3

sqlite3 db "insert into pod(node_id, name) VALUES(\"$node_id\", \"$pod_name\");"
last_pod_id=$(sqlite3 db "SELECT last_insert_rowid()")
sqlite3 db "insert into machine(docker_id, pod_id) VALUES(\"$container_docker_id\", \"$last_pod_id\");"

#redis-cli set "$line" "new-machine"