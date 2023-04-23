#!/bin/bash

read -r line

sqlite3 db "insert into machine(id, ip, status) VALUES(\"$line\", \"192.168.0.16\", 1);"

#redis-cli set "$line" "new-machine"