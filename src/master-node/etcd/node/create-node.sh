#!/bin/bash

ip_network=$1

node_name="net-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

cluster_id=$(sqlite3 db "select id from cluster limit 1")

sqlite3 db "insert into node(name, cluster_id, network, mask) values(\"$node_name\", \"$cluster_id\", \"$ip_network\", 16)"