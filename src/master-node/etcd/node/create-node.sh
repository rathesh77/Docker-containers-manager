#!/bin/bash

node_name="net-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

cluster_id=$(sqlite3 db "select id from cluster limit 1")

sqlite3 db "insert into node(name, cluster_id, network, mask) values(\"$node_name\", \"$cluster_id\", \"192.168.0.32\", 16)"