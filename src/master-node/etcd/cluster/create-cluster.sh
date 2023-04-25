#!/bin/bash

read -r line

sqlite3 db "insert into cluster(name, status) VALUES(\"$line\", 1);"

#redis-cli set "$line" "new-machine"