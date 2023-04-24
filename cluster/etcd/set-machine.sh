#!/bin/bash

read -r line

container_id=$(echo "$line" | sed -r 's/.*://')

ns=$(echo "$line" | sed -r 's/:.*//')

ns_id=$(sqlite3 db "select id from namespace where name ='$ns'";)

sqlite3 db "insert into machine(id, ip, namespace_id, status) VALUES(\"$container_id\", \"192.168.0.16\", \"$ns_id\", 1);"

#redis-cli set "$line" "new-machine"