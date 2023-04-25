#!/bin/bash

read -r line

container_id=$(echo "$line" | sed -r 's/.*://')

cluster=$(echo "$line" | sed -r 's/:.*//')

cluster_id=$(sqlite3 db "select id from cluster where name ='$cluster'";)

sqlite3 db "insert into machine(id, ip, cluster_id, status) VALUES(\"$container_id\", \"192.168.0.16\", \"$cluster_id\", 1);"

#redis-cli set "$line" "new-machine"