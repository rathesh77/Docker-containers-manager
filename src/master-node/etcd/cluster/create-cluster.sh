#!/bin/bash

read -r line

network="net-$(tr -dc A-Za-z0-9 </dev/urandom | head -c 13 ; echo '')"

docker network create "$network"

sqlite3 db "insert into cluster(name, network, status) VALUES(\"$line\",\"$network\", 1);"

#redis-cli set "$line" "new-machine"