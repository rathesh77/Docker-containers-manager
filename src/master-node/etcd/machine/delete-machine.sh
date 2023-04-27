#!/bin/bash

id=$1

sqlite3 db "delete from machine where id=\"$id\";"

#redis-cli set "$line" "new-machine"