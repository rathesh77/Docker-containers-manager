#!/bin/bash

node_id=$1
container_id=$2
container_name=$3

sqlite3 db "insert into machine(id, name, ip, node_id) VALUES(\"$container_id\", \"$container_name\", \"192.168.0.16\", \"$node_id\");"

#redis-cli set "$line" "new-machine"